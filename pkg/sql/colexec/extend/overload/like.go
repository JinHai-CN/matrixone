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

package overload

import (
	"errors"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/encoding"
	"github.com/matrixorigin/matrixone/pkg/vectorize/like"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
	"github.com/matrixorigin/matrixone/pkg/vm/register"
)

var (
	errTemp = errors.New("operator LIKE can not support for NULL now")
	errUnexpected = errors.New("unexpected case for LIKE operator")
)

func init() {
	BinOps[Like] = []*BinOp{
		{
			LeftType: types.T_char,
			RightType: types.T_char,
			ReturnType: types.T_sel,
			Fn: func(lv *vector.Vector, rv *vector.Vector, proc *process.Process, lc bool, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.(*types.Bytes), rv.Col.(*types.Bytes)
				rtl := int64(SelsType.Size)  //Type(T_sel).Length
				switch {
				// 1. []string Reg expr
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs.Lengths)) * rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs.Lengths)]
					if lv.Nsp.Any() {
						rs, err = like.SliceNullLikePure(lvs, rvs.Get(0), lv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = lv.Nsp
					} else {
						rs, err = like.SliceLikePure(lvs, rvs.Get(0), rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				// 2. string Reg expr
				case lc && rc: // in our design, this case should deal while pruning extends.
					vec, err := register.Get(proc, rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:1]
					rs, err = like.PureLikePure(lvs.Get(0), rvs.Get(0), rs)
					if err != nil {
						return nil, err
					}
					vec.SetCol(rs)
					return vec, nil
				// 3. string Reg []expr
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs.Lengths)) * rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs.Lengths)]
					if rv.Nsp.Any() {
						rs, err = like.PureLikeSliceNull(lvs.Get(0), rvs, rv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = rv.Nsp
					} else {
						rs, err = like.PureLikeSlice(lvs.Get(0), rvs, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				// 4. []string Reg []expr
				case !lc && !rc:
					vec, err := register.Get(proc, int64(len(lvs.Lengths)) * rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs.Lengths)]
					if rv.Nsp.Any() && lv.Nsp.Any() {
						nsp := lv.Nsp.Or(rv.Nsp)
						rs, err = like.SliceNullLikeSliceNull(lvs, rvs, nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = nsp
					} else if rv.Nsp.Any() && !lv.Nsp.Any() {
						rs, err = like.SliceNullLikeSliceNull(lvs, rvs, rv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = rv.Nsp
					} else if !rv.Nsp.Any() && lv.Nsp.Any() {
						rs, err = like.SliceNullLikeSliceNull(lvs, rvs, lv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = lv.Nsp
					} else {
						rs, err = like.SliceLikeSlice(lvs, rvs, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				}
				return nil, errUnexpected
			},
		},

		{
			LeftType: types.T_varchar,
			RightType: types.T_varchar,
			ReturnType: types.T_sel,
			Fn: func(lv *vector.Vector, rv *vector.Vector, proc *process.Process, lc bool, rc bool) (*vector.Vector, error) {
				lvs, rvs := lv.Col.(*types.Bytes), rv.Col.(*types.Bytes)
				rtl := int64(SelsType.Size)  //Type(T_sel).Length
				switch {
				// 1. []string Reg expr
				case !lc && rc:
					vec, err := register.Get(proc, int64(len(lvs.Lengths)) * rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(lvs.Lengths)]
					if lv.Nsp.Any() {
						rs, err = like.SliceNullLikePure(lvs, rvs.Get(0), lv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = lv.Nsp
					} else {
						rs, err = like.SliceLikePure(lvs, rvs.Get(0), rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					return vec, nil
				// 2. string Reg expr
				case lc && rc: // in our design, this case should deal while pruning extends.
					vec, err := register.Get(proc, rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:1]
					rs, err = like.PureLikePure(lvs.Get(0), rvs.Get(0), rs)
					if err != nil {
						return nil, err
					}
					vec.SetCol(rs)
					return vec, nil
				// 3. string Reg []expr
				case lc && !rc:
					vec, err := register.Get(proc, int64(len(rvs.Lengths)) * rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs.Lengths)]
					if rv.Nsp.Any() {
						rs, err = like.PureLikeSliceNull(lvs.Get(0), rvs, rv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = rv.Nsp
					} else {
						rs, err = like.PureLikeSlice(lvs.Get(0), rvs, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				// 4. []string Reg []expr
				case !lc && !rc:
					vec, err := register.Get(proc, int64(len(lvs.Lengths)) * rtl, SelsType)
					if err != nil {
						return nil, err
					}
					rs := encoding.DecodeInt64Slice(vec.Data)
					rs = rs[:len(rvs.Lengths)]
					if rv.Nsp.Any() && lv.Nsp.Any() {
						nsp := lv.Nsp.Or(rv.Nsp)
						rs, err = like.SliceNullLikeSliceNull(lvs, rvs, nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = nsp
					} else if rv.Nsp.Any() && !lv.Nsp.Any() {
						rs, err = like.SliceNullLikeSliceNull(lvs, rvs, rv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = rv.Nsp
					} else if !rv.Nsp.Any() && lv.Nsp.Any() {
						rs, err = like.SliceNullLikeSliceNull(lvs, rvs, lv.Nsp.Np, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
						vec.Nsp = lv.Nsp
					} else {
						rs, err = like.SliceLikeSlice(lvs, rvs, rs)
						if err != nil {
							return nil, err
						}
						vec.SetCol(rs)
					}
					if lv.Ref == 0 {
						register.Put(proc, lv)
					}
					if rv.Ref == 0 {
						register.Put(proc, rv)
					}
					return vec, nil
				}
				return nil, errUnexpected
			},
		},
	}
}