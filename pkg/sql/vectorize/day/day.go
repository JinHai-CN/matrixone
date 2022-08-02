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

package day

import (
	"github.com/matrixorigin/matrixone/pkg/common/container/types"
)

var (
	DateToDay     func([]types.Date, []uint8) []uint8
	DatetimeToDay func([]types.Datetime, []uint8) []uint8
)

func init() {
	DateToDay = dateToDay
	DatetimeToDay = datetimeToDay
}

func dateToDay(xs []types.Date, rs []uint8) []uint8 {
	for i, x := range xs {
		rs[i] = x.Day()
	}
	return rs
}

func datetimeToDay(xs []types.Datetime, rs []uint8) []uint8 {
	for i, x := range xs {
		rs[i] = x.Day()
	}
	return rs
}
