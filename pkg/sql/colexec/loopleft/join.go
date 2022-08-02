// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loopleft

import (
	"bytes"
	"github.com/matrixorigin/matrixone/pkg/common/container/batch"
	"github.com/matrixorigin/matrixone/pkg/common/container/vector"

	"github.com/matrixorigin/matrixone/pkg/common/hashmap"
	"github.com/matrixorigin/matrixone/pkg/sql/colexec"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
)

func String(_ any, buf *bytes.Buffer) {
	buf.WriteString(" loop left join ")
}

func Prepare(proc *process.Process, arg any) error {
	ap := arg.(*Argument)
	ap.ctr = new(container)
	ap.ctr.bat = batch.NewWithSize(len(ap.Typs))
	ap.ctr.bat.Zs = proc.GetMheap().GetSels()
	for i, typ := range ap.Typs {
		ap.ctr.bat.Vecs[i] = vector.New(typ)
	}
	return nil
}

func Call(idx int, proc *process.Process, arg any) (bool, error) {
	anal := proc.GetAnalyze(idx)
	anal.Start()
	defer anal.Stop()
	ap := arg.(*Argument)
	ctr := ap.ctr
	for {
		switch ctr.state {
		case Build:
			if err := ctr.build(ap, proc, anal); err != nil {
				ctr.state = End
				return true, err
			}
			ctr.state = Probe
		case Probe:
			bat := <-proc.Reg.MergeReceivers[0].Ch
			if bat == nil {
				ctr.state = End
				if ctr.bat != nil {
					ctr.bat.Clean(proc.Mp)
				}
				continue
			}
			if bat.Length() == 0 {
				continue
			}
			if ctr.bat.Length() == 0 {
				if err := ctr.emptyProbe(bat, ap, proc, anal); err != nil {
					ctr.state = End
					proc.SetInputBatch(nil)
					return true, err
				}

			} else {
				if err := ctr.probe(bat, ap, proc, anal); err != nil {
					ctr.state = End
					proc.SetInputBatch(nil)
					return true, err
				}
			}
			return false, nil
		default:
			proc.SetInputBatch(nil)
			return true, nil
		}
	}
}

func (ctr *container) build(ap *Argument, proc *process.Process, anal process.Analyze) error {
	var err error

	for {
		bat := <-proc.Reg.MergeReceivers[1].Ch
		if bat == nil {
			break
		}
		if bat.Length() == 0 {
			continue
		}
		anal.Input(bat)
		anal.Alloc(int64(bat.Size()))
		if ctr.bat, err = ctr.bat.Append(proc.GetMheap(), bat); err != nil {
			bat.Clean(proc.GetMheap())
			ctr.bat.Clean(proc.GetMheap())
			return err
		}
		bat.Clean(proc.GetMheap())
	}
	return nil
}

func (ctr *container) emptyProbe(bat *batch.Batch, ap *Argument, proc *process.Process, anal process.Analyze) error {
	defer bat.Clean(proc.GetMheap())
	anal.Input(bat)
	rbat := batch.NewWithSize(len(ap.Result))
	rbat.Zs = proc.GetMheap().GetSels()
	for i, rp := range ap.Result {
		if rp.Rel == 0 {
			rbat.Vecs[i] = vector.New(bat.Vecs[rp.Pos].Typ)
		} else {
			rbat.Vecs[i] = vector.New(ctr.bat.Vecs[rp.Pos].Typ)
		}
	}
	count := bat.Length()
	for i := 0; i < count; i += hashmap.UnitLimit {
		n := count - i
		if n > hashmap.UnitLimit {
			n = hashmap.UnitLimit
		}
		for k := 0; k < n; k++ {
			for j, rp := range ap.Result {
				if rp.Rel == 0 {
					if err := vector.UnionOne(rbat.Vecs[j], bat.Vecs[rp.Pos], int64(i+k), proc.GetMheap()); err != nil {
						rbat.Clean(proc.GetMheap())
						return err
					}
				} else {
					if err := vector.UnionNull(rbat.Vecs[j], nil, proc.GetMheap()); err != nil {
						rbat.Clean(proc.GetMheap())
						return err
					}
				}
			}
			rbat.Zs = append(rbat.Zs, bat.Zs[i+k])
		}
	}
	rbat.ExpandNulls()
	anal.Output(rbat)
	proc.SetInputBatch(rbat)
	return nil
}

func (ctr *container) probe(bat *batch.Batch, ap *Argument, proc *process.Process, anal process.Analyze) error {
	defer bat.Clean(proc.GetMheap())
	anal.Input(bat)
	rbat := batch.NewWithSize(len(ap.Result))
	rbat.Zs = proc.GetMheap().GetSels()
	for i, rp := range ap.Result {
		if rp.Rel == 0 {
			rbat.Vecs[i] = vector.New(bat.Vecs[rp.Pos].Typ)
		} else {
			rbat.Vecs[i] = vector.New(ctr.bat.Vecs[rp.Pos].Typ)
		}
	}
	count := bat.Length()
	for i := 0; i < count; i++ {
		flg := true
		vec, err := colexec.JoinFilterEvalExpr(bat, ctr.bat, i, proc, ap.Cond)
		if err != nil {
			return err
		}
		bs := vec.Col.([]bool)
		if len(bs) == 1 {
			if bs[0] {
				for j := 0; j < len(ctr.bat.Zs); j++ {
					flg = false
					for k, rp := range ap.Result {
						if rp.Rel == 0 {
							if err := vector.UnionOne(rbat.Vecs[k], bat.Vecs[rp.Pos], int64(i), proc.GetMheap()); err != nil {
								rbat.Clean(proc.GetMheap())
								return err
							}
						} else {
							if err := vector.UnionOne(rbat.Vecs[k], ctr.bat.Vecs[rp.Pos], int64(j), proc.GetMheap()); err != nil {
								rbat.Clean(proc.GetMheap())
								return err
							}
						}
					}
					rbat.Zs = append(rbat.Zs, ctr.bat.Zs[j])

				}
			}
		} else {
			for j, b := range bs {
				if b {
					flg = false
					for k, rp := range ap.Result {
						if rp.Rel == 0 {
							if err := vector.UnionOne(rbat.Vecs[k], bat.Vecs[rp.Pos], int64(i), proc.GetMheap()); err != nil {
								rbat.Clean(proc.GetMheap())
								return err
							}
						} else {
							if err := vector.UnionOne(rbat.Vecs[k], ctr.bat.Vecs[rp.Pos], int64(j), proc.GetMheap()); err != nil {
								rbat.Clean(proc.GetMheap())
								return err
							}
						}
					}
					rbat.Zs = append(rbat.Zs, ctr.bat.Zs[j])
				}
			}
		}
		vector.Clean(vec, proc.Mp)
		if flg {
			for k, rp := range ap.Result {
				if rp.Rel == 0 {
					if err := vector.UnionOne(rbat.Vecs[k], bat.Vecs[rp.Pos], int64(i), proc.GetMheap()); err != nil {
						rbat.Clean(proc.GetMheap())
						return err
					}
				} else {
					if err := vector.UnionNull(rbat.Vecs[k], ctr.bat.Vecs[rp.Pos], proc.GetMheap()); err != nil {
						rbat.Clean(proc.GetMheap())
						return err
					}
				}
			}
			rbat.Zs = append(rbat.Zs, bat.Zs[i])
		}
	}
	rbat.ExpandNulls()
	anal.Output(rbat)
	proc.SetInputBatch(rbat)
	return nil
}
