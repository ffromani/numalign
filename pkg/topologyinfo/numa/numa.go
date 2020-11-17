/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2020 Red Hat, Inc.
 */

package numa

import (
	"github.com/fromanirh/numalign/pkg/topologyinfo/sysfs"
)

type Nodes struct {
	Online           []int
	Possible         []int
	WithCPU          []int
	WithMemory       []int
	WithNormalMemory []int
}

func NewNodesFromSysFS(sysfsPath string) (Nodes, error) {
	sysNode := sysfs.New(sysfsPath).Join(sysfs.PathDevsSysNode)

	online, err := sysNode.ReadList("online")
	if err != nil {
		return Nodes{}, err
	}
	possible, err := sysNode.ReadList("possible")
	if err != nil {
		return Nodes{}, err
	}
	hasCPU, err := sysNode.ReadList("has_cpu")
	if err != nil {
		return Nodes{}, err
	}
	hasMemory, err := sysNode.ReadList("has_memory")
	if err != nil {
		return Nodes{}, err
	}
	hasNormalMemory, err := sysNode.ReadList("has_normal_memory")
	if err != nil {
		return Nodes{}, err
	}

	return Nodes{
		Online:           online,
		Possible:         possible,
		WithCPU:          hasCPU,
		WithMemory:       hasMemory,
		WithNormalMemory: hasNormalMemory,
	}, nil
}
