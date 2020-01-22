package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func splitCPUList(cpuList string) ([]int, error) {
	cpus, err := cpuset.Parse(cpuList)
	if err != nil {
		return nil, err
	}
	return cpus.ToSlice(), nil
}

func allowedCPUList(statusFile string) ([]int, error) {
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

func findCPUSPerNUMANode(sysfsdir string) (map[int][]int, error) {
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

func findNUMANodePerCPUs(cpusPerNUMA map[int][]int) map[int]int {
	NUMAPerCPUs := make(map[int]int)
	for nodeNum, cpus := range cpusPerNUMA {
		for _, cpu := range cpus {
			NUMAPerCPUs[cpu] = nodeNum
		}
	}
	return NUMAPerCPUs
}

func findPCIDevicesFromEnv() []string {
	var pciDevs []string
	for _, envVar := range os.Environ() {
		if !strings.HasPrefix(envVar, "PCIDEVICE_") {
			continue
		}
		pair := strings.SplitN(envVar, "=", 2)
		pciDevs = append(pciDevs, pair[1])
	}
	return pciDevs
}

func findNUMANodePerPCIDevice(sysPCIDir string, devs []string) (map[string]int, error) {
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

func checkNUMAAlignment(cpuList []int, NUMAPerCPUs map[int]int, NUMAPerDev map[string]int) (int, bool) {
	nodeNum := -1
	for _, cpuID := range cpuList {
		if nodeNum == -1 {
			nodeNum = NUMAPerCPUs[cpuID]
		} else if nodeNum != NUMAPerCPUs[cpuID] {
			return -1, false
		}
	}

	for _, devNode := range NUMAPerDev {
		// TODO: explain -1
		if devNode != -1 && nodeNum != devNode {
			return -1, false
		}
	}
	return nodeNum, true
}

func main() {
	var err error

	if _, ok := os.LookupEnv("DEBUG"); !ok {
		log.SetOutput(ioutil.Discard)
	}

	cpuIDs, err := allowedCPUList("/proc/self/status")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("CPU: allowed: %v", cpuIDs)

	cpusPerNUMA, err := findCPUSPerNUMANode("/sys/devices/system/node")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("CPU: id by NUMA node: %v", cpusPerNUMA)

	NUMAPerCPUs := findNUMANodePerCPUs(cpusPerNUMA)
	log.Printf("CPU: NUMA node by id: %v", NUMAPerCPUs)

	NUMAPerDev := make(map[string]int)
	pciDevs := findPCIDevicesFromEnv()
	if len(pciDevs) > 0 {
		log.Printf("PCI: devices: %s", strings.Join(pciDevs, " - "))

		NUMAPerDev, err = findNUMANodePerPCIDevice("/sys/bus/pci/devices/", pciDevs)
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		log.Printf("PCI: devices not specified")
	}

	nodeNum, aligned := checkNUMAAlignment(cpuIDs, NUMAPerCPUs, NUMAPerDev)
	fmt.Printf("ALIGNED=%v\n", aligned)
	fmt.Printf("NODE=%v\n", nodeNum)
	time.Sleep(24 * time.Hour)
}
