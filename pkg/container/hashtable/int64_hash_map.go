// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hashtable

import (
	"errors"
	"unsafe"
)

type Int64HashMapCell struct {
	Key    uint64
	Mapped uint64
}

type Int64HashMap struct {
	bucketCntBits uint8
	bucketCnt     uint64
	elemCnt       uint64
	maxElemCnt    uint64
	zeroCell      Int64HashMapCell
	rawData       []byte
	bucketData    []Int64HashMapCell
}

func (ht *Int64HashMap) Init() {
	const cellSize = int(unsafe.Sizeof(Int64HashMapCell{}))

	rawData := make([]byte, cellSize*kInitialBucketCnt)

	ht.bucketCntBits = kInitialBucketCntBits
	ht.bucketCnt = kInitialBucketCnt
	ht.elemCnt = 0
	ht.maxElemCnt = kInitialBucketCnt * kLoadFactorNumerator / kLoadFactorDenominator
	ht.rawData = rawData
	ht.bucketData = unsafe.Slice((*Int64HashMapCell)(unsafe.Pointer(&rawData[0])), kInitialBucketCnt)
}

func (ht *Int64HashMap) Insert(hash uint64, keyPtr unsafe.Pointer) (inserted bool, value uint64) {
	key := *(*uint64)(keyPtr)
	if key == 0 {
		if ht.zeroCell.Mapped == 0 {
			ht.elemCnt++
			ht.zeroCell.Mapped = ht.elemCnt
		}
		value = ht.zeroCell.Mapped
		return
	}

	ht.resizeOnDemand(1)

	if hash == 0 {
		hash = crc32Int64Hash(key)
	}

	inserted, _, cell := ht.findBucket(hash, key)
	if inserted {
		ht.elemCnt++
		cell.Key = key
		cell.Mapped = ht.elemCnt
	}

	value = cell.Mapped

	return
}

func (ht *Int64HashMap) InsertBatch(n int, hashes []uint64, keysPtr unsafe.Pointer, inserted []uint8, values []uint64) {
	ht.resizeOnDemand(n)

	if hashes[0] == 0 {
		crc32Int64BatchHash(keysPtr, &hashes[0], n)
	}

	keys := unsafe.Slice((*uint64)(keysPtr), n)

	for i, key := range keys {
		if key == 0 {
			if ht.zeroCell.Mapped == 0 {
				ht.elemCnt++
				inserted[i] = 1
				ht.zeroCell.Mapped = ht.elemCnt
			}
			values[i] = ht.zeroCell.Mapped
		} else {
			isInserted, _, cell := ht.findBucket(hashes[i], key)
			if isInserted {
				ht.elemCnt++
				inserted[i] = 1
				cell.Key = key
				cell.Mapped = ht.elemCnt
			}
			values[i] = cell.Mapped
		}
	}
}

func (ht *Int64HashMap) Find(hash uint64, keyPtr unsafe.Pointer) uint64 {
	key := *(*uint64)(keyPtr)
	if key == 0 {
		return ht.zeroCell.Mapped
	}

	if hash == 0 {
		hash = crc32Int64Hash(key)
	}

	_, _, cell := ht.findBucket(hash, key)
	return cell.Mapped
}

func (ht *Int64HashMap) FindBatch(n int, hashes []uint64, keysPtr unsafe.Pointer, values []uint64) {
	if hashes[0] == 0 {
		crc32Int64BatchHash(keysPtr, &hashes[0], n)
	}

	keys := unsafe.Slice((*uint64)(keysPtr), n)
	for i, key := range keys {
		if key == 0 {
			values[i] = ht.zeroCell.Mapped
		} else {
			_, _, cell := ht.findBucket(hashes[i], key)
			values[i] = cell.Mapped
		}
	}
}

func (ht *Int64HashMap) findBucket(hash uint64, key uint64) (empty bool, idx uint64, cell *Int64HashMapCell) {
	mask := ht.bucketCnt - 1
	var equal bool
	for idx = hash & mask; true; idx = (idx + 1) & mask {
		cell = &ht.bucketData[idx]
		empty, equal = cell.Key == 0, cell.Key == key
		if empty || equal {
			return
		}
	}

	return
}

func (ht *Int64HashMap) resizeOnDemand(n int) {
	targetCnt := ht.elemCnt + uint64(n)
	if targetCnt <= ht.maxElemCnt {
		return
	}

	newBucketCntBits := ht.bucketCntBits + 2
	newBucketCnt := uint64(1) << newBucketCntBits
	newMaxElemCnt := newBucketCnt * kLoadFactorNumerator / kLoadFactorDenominator
	for newMaxElemCnt < targetCnt {
		newBucketCntBits++
		newBucketCnt <<= 1
		newMaxElemCnt = newBucketCnt * kLoadFactorNumerator / kLoadFactorDenominator
	}

	const cellSize = int(unsafe.Sizeof(Int64HashMapCell{}))

	oldBucketCnt := ht.bucketCnt
	oldBucketData := ht.bucketData

	newRawData := make([]byte, uint64(cellSize)*newBucketCnt)
	newBucketData := unsafe.Slice((*Int64HashMapCell)(unsafe.Pointer(&newRawData[0])), newBucketCnt)

	ht.bucketCntBits = newBucketCntBits
	ht.bucketCnt = newBucketCnt
	ht.maxElemCnt = newMaxElemCnt
	ht.rawData = newRawData
	ht.bucketData = newBucketData

	var hashes [256]uint64

	var i uint64
	for i = 0; i < oldBucketCnt; i += 256 {
		cells := oldBucketData[i : i+256]
		crc32Int64CellBatchHash(unsafe.Pointer(&cells[0]), &hashes[0], 256)
		for j := range cells {
			cell := &cells[j]
			if cell.Key != 0 {
				_, newIdx, _ := ht.findBucket(hashes[j], cell.Key)
				ht.bucketData[newIdx] = *cell
			}
		}
	}
}

func (ht *Int64HashMap) Cardinality() uint64 {
	return ht.elemCnt
}

type Int64HashMapIterator struct {
	table   *Int64HashMap
	pos     uint64
	visited bool
}

func (it *Int64HashMapIterator) Init(ht *Int64HashMap) {
	it.table = ht
}

func (it *Int64HashMapIterator) Next() (cell *Int64HashMapCell, err error) {
	if !it.visited {
		it.visited = true
		if it.table.zeroCell.Mapped > 0 {
			cell = &it.table.zeroCell
			return
		}
	}

	for it.pos < it.table.bucketCnt {
		cell = &it.table.bucketData[it.pos]
		if cell.Key != 0 {
			break
		}
		it.pos++
	}

	if it.pos >= it.table.bucketCnt {
		err = errors.New("out of range")
		return
	}

	it.pos++

	return
}
