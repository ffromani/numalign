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

package cpus

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fromanirh/cpuset"
)

/*
 * keep this handy:
 * https://www.kernel.org/doc/html/latest/admin-guide/cputopology.html
 */
const (
	PathDevsSysCPU  = "/sys/devices/system/cpu"
	PathDevsSysNode = "/sys/devices/system/node"
)

type CPUIdList []int

type CPUs struct {
	Present      CPUIdList
	Online       CPUIdList
	CoreCPUs     map[int]CPUIdList // aka thread_siblings
	PackageCPUs  map[int]CPUIdList // aka core_siblings
	Packages     int
	NUMANodes    int
	NUMANodeCPUs map[int]CPUIdList
}

func NewCPUs(sysfs string) (*CPUs, error) {
	sysfsCPUPath := filepath.Join(sysfs, PathDevsSysCPU)
	present, err := readCPUList(filepath.Join(sysfsCPUPath, "present"))
	if err != nil {
		return nil, err
	}
	online, err := readCPUList(filepath.Join(sysfsCPUPath, "online"))
	if err != nil {
		return nil, err
	}

	sysfsNodePath := filepath.Join(sysfs, PathDevsSysNode)
	nodes, err := countNUMANodes(sysfsNodePath)
	if err != nil {
		return nil, err
	}

	packages := make(map[string]CPUIdList)
	coreCPUs := make(map[int]CPUIdList)
	packageCPUs := make(map[int]CPUIdList)
	for _, cpuId := range online {
		sysfsCPUIdPath := pathSysCPUxTopology(sysfsCPUPath, cpuId)
		cpuThreads, err := readCPUList(filepath.Join(sysfsCPUIdPath, "thread_siblings_list"))
		if err != nil {
			return nil, err
		}
		cpuCores, err := readCPUList(filepath.Join(sysfsCPUIdPath, "core_siblings_list"))
		if err != nil {
			return nil, err
		}
		physPackageId, err := readSysFSFile(filepath.Join(sysfsCPUIdPath, "physical_package_id"))
		if err != nil {
			return nil, err
		}

		coreCPUs[cpuId] = cpuThreads
		packageCPUs[cpuId] = cpuCores

		cpusPerPhysPkg := packages[physPackageId]
		cpusPerPhysPkg = append(cpusPerPhysPkg, cpuId)
		packages[physPackageId] = cpusPerPhysPkg
	}

	numaNodeCPUs := make(map[int]CPUIdList)
	for node := 0; node < nodes; node++ {
		cpus, err := readCPUList(filepath.Join(pathSysNodex(sysfsNodePath, node), "cpulist"))
		if err != nil {
			return nil, err
		}
		numaNodeCPUs[node] = cpus
	}

	return &CPUs{
		Present:      present,
		Online:       online,
		CoreCPUs:     coreCPUs,
		PackageCPUs:  packageCPUs,
		Packages:     len(packages),
		NUMANodes:    nodes,
		NUMANodeCPUs: numaNodeCPUs,
	}, nil
}

func readSysFSFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func readCPUList(path string) (CPUIdList, error) {
	data, err := readSysFSFile(path)
	if err != nil {
		return nil, err
	}
	cpus, err := cpuset.Parse(data)
	if err != nil {
		return nil, err
	}
	return CPUIdList(cpus), nil
}

func pathSysCPUxTopology(sysfsCPUPath string, cpuId int) string {
	return filepath.Join(sysfsCPUPath, fmt.Sprintf("cpu%d", cpuId), "topology")
}

func pathSysNodex(sysfsNodePath string, nodeId int) string {
	return filepath.Join(sysfsNodePath, fmt.Sprintf("node%d", nodeId))
}

func countNUMANodes(nodepath string) (int, error) {
	nodes := 0
	entries, err := ioutil.ReadDir(nodepath)
	if err != nil {
		return nodes, err
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "node") {
			nodes++
		}
	}
	return nodes, nil
}
