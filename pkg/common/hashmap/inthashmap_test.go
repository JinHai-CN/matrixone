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

package hashmap

import (
	"github.com/matrixorigin/matrixone/pkg/common/container/types"
	"github.com/matrixorigin/matrixone/pkg/common/container/vector"
	"testing"

	"github.com/matrixorigin/matrixone/pkg/vm/mheap"
	"github.com/matrixorigin/matrixone/pkg/vm/mmu/guest"
	"github.com/matrixorigin/matrixone/pkg/vm/mmu/host"
	"github.com/stretchr/testify/require"
)

func TestIntHashMap_Iterator(t *testing.T) {
	{
		mp := NewIntHashMap(false)
		hm := host.New(1 << 30)
		gm := guest.New(1<<30, hm)
		m := mheap.New(gm)
		rowCount := 10
		vecs := []*vector.Vector{
			newVector(rowCount, types.T_int32.ToType(), m, false, []int32{
				-1, -1, -1, 2, 2, 2, 3, 3, 3, 4,
			}),
			newVector(rowCount, types.T_uint32.ToType(), m, false, []uint32{
				1, 1, 1, 2, 2, 2, 3, 3, 3, 4,
			}),
		}
		itr := mp.NewIterator(0, 0)
		vs, _ := itr.Insert(0, rowCount, vecs)
		require.Equal(t, []uint64{1, 1, 1, 2, 2, 2, 3, 3, 3, 4}, vs)
		vs, _ = itr.Find(0, rowCount, vecs, nil)
		require.Equal(t, []uint64{1, 1, 1, 2, 2, 2, 3, 3, 3, 4}, vs)
		for _, vec := range vecs {
			vec.Free(m)
		}
		require.Equal(t, int64(0), m.Size())
	}
	{
		mp := NewIntHashMap(true)
		ts := []types.Type{
			types.New(types.T_int8, 0, 0, 0),
			types.New(types.T_int16, 0, 0, 0),
		}
		hm := host.New(1 << 30)
		gm := guest.New(1<<30, hm)
		m := mheap.New(gm)
		vecs := newVectorsWithNull(ts, false, Rows, m)
		itr := mp.NewIterator(0, 0)
		vs, _ := itr.Insert(0, Rows, vecs)
		require.Equal(t, []uint64{1, 2, 1, 3, 1, 4, 1, 5, 1, 6}, vs[:Rows])
		vs, _ = itr.Find(0, Rows, vecs, nil)
		require.Equal(t, []uint64{1, 2, 1, 3, 1, 4, 1, 5, 1, 6}, vs[:Rows])
		for _, vec := range vecs {
			vec.Free(m)
		}
		require.Equal(t, int64(0), m.Size())
	}
	{
		mp := NewIntHashMap(true)
		ts := []types.Type{
			types.New(types.T_int64, 0, 0, 0),
		}
		hm := host.New(1 << 30)
		gm := guest.New(1<<30, hm)
		m := mheap.New(gm)
		vecs := newVectorsWithNull(ts, false, Rows, m)
		itr := mp.NewIterator(0, 0)
		vs, _ := itr.Insert(0, Rows, vecs)
		require.Equal(t, []uint64{1, 2, 1, 3, 1, 4, 1, 5, 1, 6}, vs[:Rows])
		vs, _ = itr.Find(0, Rows, vecs, nil)
		require.Equal(t, []uint64{1, 2, 1, 3, 1, 4, 1, 5, 1, 6}, vs[:Rows])
		for _, vec := range vecs {
			vec.Free(m)
		}
		require.Equal(t, int64(0), m.Size())
	}
	{
		mp := NewIntHashMap(true)
		ts := []types.Type{
			types.New(types.T_char, 1, 0, 0),
		}
		hm := host.New(1 << 30)
		gm := guest.New(1<<30, hm)
		m := mheap.New(gm)
		vecs := newVectorsWithNull(ts, false, Rows, m)
		itr := mp.NewIterator(0, 0)
		vs, _ := itr.Insert(0, Rows, vecs)
		require.Equal(t, []uint64{1, 2, 1, 3, 1, 4, 1, 5, 1, 6}, vs[:Rows])
		vs, _ = itr.Find(0, Rows, vecs, nil)
		require.Equal(t, []uint64{1, 2, 1, 3, 1, 4, 1, 5, 1, 6}, vs[:Rows])
		for _, vec := range vecs {
			vec.Free(m)
		}
		require.Equal(t, int64(0), m.Size())
	}
	{
		mp := NewIntHashMap(true)
		ts := []types.Type{
			types.New(types.T_char, 1, 0, 0),
		}
		hm := host.New(1 << 30)
		gm := guest.New(1<<30, hm)
		m := mheap.New(gm)
		vecs := newVectors(ts, false, Rows, m)
		itr := mp.NewIterator(0, 0)
		vs, _ := itr.Insert(0, Rows, vecs)
		require.Equal(t, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, vs[:Rows])
		vs, _ = itr.Find(0, Rows, vecs, nil)
		require.Equal(t, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, vs[:Rows])
		for _, vec := range vecs {
			vec.Free(m)
		}
		require.Equal(t, int64(0), m.Size())
	}
	{
		mp := NewIntHashMap(false)
		ts := []types.Type{
			types.New(types.T_char, 1, 0, 0),
		}
		hm := host.New(1 << 30)
		gm := guest.New(1<<30, hm)
		m := mheap.New(gm)
		vecs := newVectorsWithNull(ts, false, Rows, m)
		itr := mp.NewIterator(0, 0)
		vs, _ := itr.Insert(0, Rows, vecs)
		require.Equal(t, []uint64{0, 1, 0, 2, 0, 3, 0, 4, 0, 5}, vs[:Rows])
		vs, _ = itr.Find(0, Rows, vecs, nil)
		require.Equal(t, []uint64{0, 1, 0, 2, 0, 3, 0, 4, 0, 5}, vs[:Rows])
		for _, vec := range vecs {
			vec.Free(m)
		}
		require.Equal(t, int64(0), m.Size())
	}
}
