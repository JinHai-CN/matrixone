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
	"github.com/matrixorigin/matrixone/pkg/storage"
	"github.com/matrixorigin/matrixone/pkg/vm/mempool"
	"github.com/matrixorigin/matrixone/pkg/vm/mmu/host"
)

var GlobalSystemVariables SystemVariables

// HostMmu host memory
var HostMmu *host.Mmu = nil

// Mempool memory pool
var Mempool *mempool.Mempool = nil

// StorageEngine Storage Engine
var StorageEngine storage.Engine

// ClusterNodes Cluster Nodes
var ClusterNodes storage.Nodes

type ParameterUnit struct {
	SV *SystemVariables

	//host memory
	HostMmu *host.Mmu

	//mempool
	Mempool *mempool.Mempool

	//Storage Engine
	StorageEngine storage.Engine

	//Cluster Nodes
	ClusterNodes storage.Nodes
}

func NewParameterUnit(sv *SystemVariables, hostMmu *host.Mmu, mempool *mempool.Mempool, storageEngine storage.Engine, clusterNodes storage.Nodes) *ParameterUnit {
	return &ParameterUnit{
		SV:            sv,
		HostMmu:       hostMmu,
		Mempool:       mempool,
		StorageEngine: storageEngine,
		ClusterNodes:  clusterNodes,
	}
}
