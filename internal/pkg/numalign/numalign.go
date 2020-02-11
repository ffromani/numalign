package numalign

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
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
	return cpus.ToSlice(), nil
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

func GetCPUToNUMANodeMap(sysNodeDir string) (map[int]int, error) {
	cpusPerNUMA, err := GetCPUsPerNUMANode(sysNodeDir)
	if err != nil {
		return nil, err
	}
	CPUToNUMANode := GetCPUNUMANodes(cpusPerNUMA)
	return CPUToNUMANode, nil
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
		return nil, fmt.Errorf("No PCI devices detected")
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
		cpusPerNUMA[nodeNum] = cpus.ToSlice()
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

func Execute() {
	var err error

	if _, ok := os.LookupEnv("NUMALIGN_DEBUG"); !ok {
		log.SetOutput(ioutil.Discard)
	}

	cpuIDs, err := GetAllowedCPUList("/proc/self/status")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("CPU: allowed: %v", cpuIDs)

	CPUToNUMANode, err := GetCPUToNUMANodeMap("/sys/devices/system/node")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("CPU: NUMA node by id: %v", CPUToNUMANode)

	pciDevs := GetPCIDevicesFromEnv(os.Environ())
	NUMAPerDev, err := GetPCIDeviceToNumaNodeMap("/sys/bus/pci/devices/", pciDevs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	R := Resources{
		CPUToNUMANode:     CPUToNUMANode,
		PCIDevsToNUMANode: NUMAPerDev,
	}
	nodeNum, aligned := R.CheckAlignment()
	fmt.Printf("STATUS ALIGNED=%v\n", aligned)
	if aligned {
		fmt.Printf("NUMA NODE=%v\n", nodeNum)
	} else {
		fmt.Printf("%s", R.String())
		os.Exit(99)
	}
}
