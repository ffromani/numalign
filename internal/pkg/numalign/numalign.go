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

package numalign

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/fromanirh/cpuset"
)

const (
	ProcStatusFile          = "/proc/self/status"
	SysDevicesSystemNodeDir = "/sys/devices/system/node"
	SysBusPCIDevicesDir     = "/sys/bus/pci/devices/"
)

func splitCPUList(cpuList string) ([]int, error) {
	cpus, err := cpuset.Parse(cpuList)
	if err != nil {
		return nil, err
	}
	return cpus, nil
}

type Resources struct {
	CPUToNUMANode     map[int]int
	PCIDevsToNUMANode map[string]int
}

func (R *Resources) CheckAlignment() (int, bool) {
	nodeNum := -1
	for _, cpuNode := range R.CPUToNUMANode {
		if nodeNum == -1 {
			nodeNum = cpuNode
		} else if nodeNum != cpuNode {
			return -1, false
		}
	}
	for _, devNode := range R.PCIDevsToNUMANode {
		// TODO: explain -1
		if devNode != -1 && nodeNum != devNode {
			return -1, false
		}
	}
	return nodeNum, true
}

func (R *Resources) String() string {
	var b strings.Builder
	// To store the keys in slice in sorted order
	var cpuKeys []int
	for ck := range R.CPUToNUMANode {
		cpuKeys = append(cpuKeys, ck)
	}
	sort.Ints(cpuKeys)
	for _, k := range cpuKeys {
		nodeNum := R.CPUToNUMANode[k]
		b.WriteString(fmt.Sprintf("CPU cpu#%03d=%02d\n", k, nodeNum))
	}
	var pciKeys []string
	for pk := range R.PCIDevsToNUMANode {
		pciKeys = append(pciKeys, pk)
	}
	sort.Strings(pciKeys)
	for _, k := range pciKeys {
		nodeNum := R.PCIDevsToNUMANode[k]
		b.WriteString(fmt.Sprintf("PCI %s=%02d\n", k, nodeNum))
	}
	return b.String()
}

func GetAllowedCPUList(statusFile string) ([]int, error) {
	var cpuIDs []int
	var err error
	content, err := ioutil.ReadFile(statusFile)
	if err != nil {
		return cpuIDs, err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Cpus_allowed_list") {
			pair := strings.SplitN(line, ":", 2)
			return splitCPUList(strings.TrimSpace(pair[1]))
		}
	}
	return cpuIDs, fmt.Errorf("malformed status file: %s", statusFile)
}

func GetCPUToNUMANodeMap(sysNodeDir string, cpuIDs []int) (map[int]int, error) {
	cpusPerNUMA, err := GetCPUsPerNUMANode(sysNodeDir)
	if err != nil {
		return nil, err
	}
	CPUToNUMANode := GetCPUNUMANodes(cpusPerNUMA)

	// filter out only the allowed CPUs
	CPUMap := make(map[int]int)
	for _, cpuID := range cpuIDs {
		_, ok := CPUToNUMANode[cpuID]
		if !ok {
			return nil, fmt.Errorf("CPU %d not found on NUMA map: %v", cpuID, CPUToNUMANode)
		}
		CPUMap[cpuID] = CPUToNUMANode[cpuID]
	}
	return CPUMap, nil
}

func GetPCIDevicesFromEnv(environ []string) []string {
	var pciDevs []string
	for _, envVar := range environ {
		if !strings.HasPrefix(envVar, "PCIDEVICE_") {
			continue
		}
		pair := strings.SplitN(envVar, "=", 2)
		pciDevs = append(pciDevs, pair[1])
	}
	return pciDevs
}

func GetPCIDeviceToNumaNodeMap(sysBusPCIDir string, pciDevs []string) (map[string]int, error) {
	if len(pciDevs) == 0 {
		log.Printf("PCI: devices: none found - SKIP")
		return make(map[string]int), nil
	}
	log.Printf("PCI: devices: %s", strings.Join(pciDevs, " - "))

	NUMAPerDev, err := GetPCIDeviceNUMANode(sysBusPCIDir, pciDevs)
	if err != nil {
		return nil, err
	}
	return NUMAPerDev, nil
}

func GetPCIDeviceNUMANode(sysPCIDir string, devs []string) (map[string]int, error) {
	NUMAPerDev := make(map[string]int)
	for _, dev := range devs {
		content, err := ioutil.ReadFile(filepath.Join(sysPCIDir, dev, "numa_node"))
		if err != nil {
			return nil, err
		}
		nodeNum, err := strconv.Atoi(strings.TrimSpace(string(content)))
		if err != nil {
			return nil, err
		}
		NUMAPerDev[dev] = nodeNum
	}
	return NUMAPerDev, nil
}

func GetCPUsPerNUMANode(sysfsdir string) (map[int][]int, error) {
	pattern := filepath.Join(sysfsdir, "node*")
	nodes, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	cpusPerNUMA := make(map[int][]int)
	for _, node := range nodes {
		_, nodeID := filepath.Split(node)
		nodeNum, err := strconv.Atoi(strings.TrimSpace(nodeID[4:]))
		if err != nil {
			return nil, err
		}
		cpuList := filepath.Join(node, "cpulist")
		content, err := ioutil.ReadFile(cpuList)
		if err != nil {
			return nil, err
		}
		cpus, err := cpuset.Parse(strings.TrimSpace(string(content)))
		if err != nil {
			return nil, err
		}
		cpusPerNUMA[nodeNum] = cpus
	}
	return cpusPerNUMA, nil
}

func GetCPUNUMANodes(cpusPerNUMA map[int][]int) map[int]int {
	CPUToNUMANode := make(map[int]int)
	for nodeNum, cpus := range cpusPerNUMA {
		for _, cpu := range cpus {
			CPUToNUMANode[cpu] = nodeNum
		}
	}
	return CPUToNUMANode
}

func Execute() int {
	var err error

	var pidStrings []string
	if len(os.Args) > 1 {
		pidStrings = append(pidStrings, os.Args[1:]...)
	} else {
		pidStrings = append(pidStrings, "self")
	}

	var refCpuIDs []int
	refCpuIDs, err = GetAllowedCPUList(filepath.Join("/proc", pidStrings[0], "status"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("CPU: allowed for %q: %v", pidStrings[0], refCpuIDs)

	for _, pidString := range pidStrings[1:] {
		cpuIDs, err := GetAllowedCPUList(filepath.Join("/proc", pidString, "status"))
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Printf("CPU: allowed for %q: %v", pidString, cpuIDs)

		if !reflect.DeepEqual(refCpuIDs, cpuIDs) {
			log.Fatalf("CPU: allowed set differs pid %q (%v) pid %q (%v)", pidStrings[0], refCpuIDs, pidString, cpuIDs)
		}
	}

	CPUToNUMANode, err := GetCPUToNUMANodeMap(SysDevicesSystemNodeDir, refCpuIDs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("CPU: NUMA node by id: %v", CPUToNUMANode)

	pciDevs := GetPCIDevicesFromEnv(os.Environ())
	NUMAPerDev, err := GetPCIDeviceToNumaNodeMap(SysBusPCIDevicesDir, pciDevs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	R := Resources{
		CPUToNUMANode:     CPUToNUMANode,
		PCIDevsToNUMANode: NUMAPerDev,
	}
	nodeNum, aligned := R.CheckAlignment()
	fmt.Printf("STATUS ALIGNED=%v\n", aligned)
	if !aligned {
		fmt.Printf("%s", R.String())
		return 1
	}

	fmt.Printf("NUMA NODE=%v\n", nodeNum)
	return 0
}
