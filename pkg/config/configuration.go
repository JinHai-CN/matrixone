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

package config

import (
	"github.com/matrixorigin/matrixone/pkg/catalog"
	"github.com/matrixorigin/matrixone/pkg/vm/engine"
	"github.com/matrixorigin/matrixone/pkg/vm/mempool"
	"github.com/matrixorigin/matrixone/pkg/vm/metadata"
	"github.com/matrixorigin/matrixone/pkg/vm/mmu/host"
)

var GlobalSystemVariables SystemVariables

//host memory
var HostMmu *host.Mmu = nil

//mempool
var Mempool  *mempool.Mempool = nil

//Storage Engine
var StorageEngine engine.Engine

//Cluster Nodes
var ClusterNodes metadata.Nodes

//cube catalog
var ClusterCatalog *catalog.Catalog

/**
check if x in a slice
*/
func isInSlice(x string, arr []string) bool {
	for _, y := range arr {
		if x == y {
			return true
		}
	}
	return false
}

/**
check if x in a slice
*/
func isInSliceBool(x bool, arr []bool) bool {
	for _, y := range arr {
		if x == y {
			return true
		}
	}
	return false
}

/**
check if x in a slice
*/
func isInSliceInt64(x int64, arr []int64) bool {
	for _, y := range arr {
		if x == y {
			return true
		}
	}
	return false
}

type ParameterUnit struct {
	SV *SystemVariables

	//host memory
	HostMmu *host.Mmu

	//mempool
	Mempool  *mempool.Mempool

	//Storage Engine
	StorageEngine engine.Engine

	//Cluster Nodes
	ClusterNodes metadata.Nodes

	//Cube Catalog
	ClusterCatalog *catalog.Catalog
}

func NewParameterUnit(sv *SystemVariables, hostMmu *host.Mmu, mempool *mempool.Mempool, storageEngine engine.Engine, clusterNodes metadata.Nodes, catalogRef *catalog.Catalog) *ParameterUnit {
	return &ParameterUnit{
		SV:            sv,
		HostMmu:      hostMmu,
		Mempool: mempool,
		StorageEngine: storageEngine,
		ClusterNodes:  clusterNodes,
		ClusterCatalog: catalogRef,
	}
}