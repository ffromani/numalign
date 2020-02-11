package numalign

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestResourceString(t *testing.T) {
	var r1 Resources
	s1 := r1.String()
	if s1 != "" {
		t.Errorf("empty Resource has unexpected output: %s", s1)
	}

	r2 := Resources{
		CPUToNUMANode: map[int]int{
			0: 0,
			1: 0,
		},
		PCIDevsToNUMANode: map[string]int{
			"3c:00.0": 0,
		},
	}
	s2 := r2.String()
	expected := `CPU cpu#000=00
CPU cpu#001=00
PCI 3c:00.0=00
`
	if s2 != expected {
		t.Errorf("initialzed Resource has unexpected output: %s (expected %s)", s2, expected)
	}

}

func TestResourceAlignment(t *testing.T) {
	r0 := Resources{}
	if nodeNum, aligned := r0.CheckAlignment(); nodeNum != -1 || !aligned {
		t.Errorf("empty resources should be considered aligned")
	}

	// all aligned
	r1 := Resources{
		CPUToNUMANode: map[int]int{
			0: 0,
			1: 0,
			2: 0,
			3: 0,
		},
		PCIDevsToNUMANode: map[string]int{
			"3c:00.0": 0,
			"00:1f.6": 0,
		},
	}
	node, aligned := r1.CheckAlignment()
	if !aligned {
		t.Errorf("aligned resourced misdetected unaligned")
	}
	if node != 0 {
		t.Errorf("resources aligned on unexpected node %d (should be 0)", node)
	}

	// cpu cores misaligned
	r2 := Resources{
		CPUToNUMANode: map[int]int{
			0: 0,
			1: 0,
			2: 0,
			3: 2,
		},
		PCIDevsToNUMANode: map[string]int{
			"3c:00.0": 0,
			"00:1f.6": 0,
		},
	}
	// we don't care about the node on unaligned resources
	_, aligned = r2.CheckAlignment()
	if aligned {
		t.Errorf("unaligned CPU cores misdetected aligned")
	}

	// PCI device misaligned
	r3 := Resources{
		CPUToNUMANode: map[int]int{
			0: 0,
			1: 0,
			2: 0,
			3: 0,
		},
		PCIDevsToNUMANode: map[string]int{
			"3c:00.0": 2,
			"00:1f.6": 0,
		},
	}
	// we don't care about the node on unaligned resources
	_, aligned = r3.CheckAlignment()
	if aligned {
		t.Errorf("unaligned PCI device misdetected aligned")
	}

	// CPU core AND PCI device misaligned
	r4 := Resources{
		CPUToNUMANode: map[int]int{
			0: 0,
			1: 0,
			2: 1,
			3: 0,
		},
		PCIDevsToNUMANode: map[string]int{
			"3c:00.0": 0,
			"00:1f.6": 1,
		},
	}
	// we don't care about the node on unaligned resources
	_, aligned = r4.CheckAlignment()
	if aligned {
		t.Errorf("unaligned CPU core AND PCI device misdetected aligned")
	}
}

func TestGetAllowedCPUList(t *testing.T) {
	cpus, err := GetAllowedCPUList("/error/proc/not/mounted")
	if err == nil {
		t.Errorf("unexpected success reading inexistent file!")
	}
	if len(cpus) != 0 {
		t.Errorf("misdetected detected allowed CPU list from inexistent file")
	}

	cpus, err = GetAllowedCPUList("/proc/cpuinfo")
	if err == nil {
		t.Errorf("malformed status file misdetected from malformed file")
	}
	if len(cpus) != 0 {
		t.Errorf("misdetected detected allowed CPU list")
	}

	cpus, err = GetAllowedCPUList("/proc/self/status")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(cpus) < 1 {
		t.Errorf("not detected allowed CPU list")
	}
}

func TestCPUToNUMANodes(t *testing.T) {
	r0 := GetCPUNUMANodes(map[int][]int{})
	expected0 := map[int]int{}
	if !reflect.DeepEqual(r0, expected0) {
		t.Errorf("maps are different: got %v expected %v", r0, expected0)
	}

	r1 := GetCPUNUMANodes(map[int][]int{
		0: []int{0, 1, 2, 3},
		1: []int{4, 5, 6, 7},
	})
	expected1 := map[int]int{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
		4: 1,
		5: 1,
		6: 1,
		7: 1,
	}
	if !reflect.DeepEqual(r1, expected1) {
		t.Errorf("maps are different: got %v expected %v", r1, expected1)
	}

	// offline CPUs?
	r2 := GetCPUNUMANodes(map[int][]int{
		0: []int{0, 1, 2, 3},
		1: []int{6, 7},
	})
	expected2 := map[int]int{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
		6: 1,
		7: 1,
	}
	if !reflect.DeepEqual(r2, expected2) {
		t.Errorf("maps are different: got %v expected %v", r2, expected2)
	}

	r3 := GetCPUNUMANodes(map[int][]int{
		0: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	})
	expected3 := map[int]int{
		0:  0,
		1:  0,
		2:  0,
		3:  0,
		4:  0,
		5:  0,
		6:  0,
		7:  0,
		8:  0,
		9:  0,
		10: 0,
		11: 0,
		12: 0,
		13: 0,
		14: 0,
		15: 0,
	}
	if !reflect.DeepEqual(r3, expected3) {
		t.Errorf("maps are different: got %v expected %v", r3, expected3)
	}

	r4 := GetCPUNUMANodes(map[int][]int{
		0:  []int{0},
		1:  []int{1},
		2:  []int{2},
		3:  []int{3},
		4:  []int{4},
		5:  []int{5},
		6:  []int{6},
		7:  []int{7},
		8:  []int{8},
		9:  []int{9},
		10: []int{10},
		11: []int{11},
		12: []int{12},
		13: []int{13},
		14: []int{14},
		15: []int{15},
	})
	expected4 := map[int]int{
		0:  0,
		1:  1,
		2:  2,
		3:  3,
		4:  4,
		5:  5,
		6:  6,
		7:  7,
		8:  8,
		9:  9,
		10: 10,
		11: 11,
		12: 12,
		13: 13,
		14: 14,
		15: 15,
	}
	if !reflect.DeepEqual(r4, expected4) {
		t.Errorf("maps are different: got %v expected %v", r4, expected4)
	}
}

func TestCPUToNumaNode(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	basedir, _ := filepath.Split(filename)

	testDataDir := filepath.Join(basedir, "..", "test", "data")

	fakeSingleSysNodeDir := filepath.Join(testDataDir, "single", SysDevicesSystemNodeDir)
	CPUToNUMANode, err := GetCPUToNUMANodeMap(fakeSingleSysNodeDir)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedSingle := map[int]int{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
		6: 0,
		7: 0,
	}
	if !reflect.DeepEqual(CPUToNUMANode, expectedSingle) {
		t.Errorf("maps are different: got %v expected %v", CPUToNUMANode, expectedSingle)
	}

	fakeMultiSysNodeDir := filepath.Join(testDataDir, "multi", SysDevicesSystemNodeDir)
	CPUToNUMANode, err = GetCPUToNUMANodeMap(fakeMultiSysNodeDir)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedMulti := map[int]int{
		0:  0,
		1:  0,
		2:  0,
		3:  0,
		4:  0,
		5:  0,
		6:  0,
		7:  0,
		8:  1,
		9:  1,
		10: 1,
		11: 1,
		12: 1,
		13: 1,
		14: 1,
		15: 1,
	}
	if !reflect.DeepEqual(CPUToNUMANode, expectedMulti) {
		t.Errorf("maps are different: got %v expected %v", CPUToNUMANode, expectedMulti)
	}
}

func TestGetPCIDevsFromEnv(t *testing.T) {
	var devs []string

	devs = GetPCIDevicesFromEnv([]string{})
	if len(devs) > 0 {
		t.Errorf("unexpected devices: %v", devs)
	}

	devs = GetPCIDevicesFromEnv([]string{"PATH=/bin:/sbin", "FOO=bar"})
	if len(devs) > 0 {
		t.Errorf("unexpected devices: %v", devs)
	}

	devs = GetPCIDevicesFromEnv([]string{"PATH=/bin:/sbin", "FOO=bar", "PCIDEVICE_FOO=0000:00:00.0"})
	if len(devs) != 1 && devs[0] != "00:00.0" {
		t.Errorf("unexpected devices: %v", devs)
	}
}

func TestGetPCIDeviceNUMANode(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	basedir, _ := filepath.Split(filename)

	testDataDir := filepath.Join(basedir, "..", "test", "data")

	fakeSysPCIDir := filepath.Join(testDataDir, SysBusPCIDevicesDir)

	NUMAPerDev1, err := GetPCIDeviceToNumaNodeMap(fakeSysPCIDir, []string{"0000:05:00.0"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected1 := map[string]int{
		"0000:05:00.0": -1,
	}
	if !reflect.DeepEqual(NUMAPerDev1, expected1) {
		t.Errorf("maps are different: got %v expected %v", NUMAPerDev1, expected1)
	}

	NUMAPerDev2, err := GetPCIDeviceToNumaNodeMap(fakeSysPCIDir, []string{"1000:00:1f.3", "1000:3c:00.0"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected2 := map[string]int{
		"1000:00:1f.3": 1,
		"1000:3c:00.0": 1,
	}
	if !reflect.DeepEqual(NUMAPerDev2, expected2) {
		t.Errorf("maps are different: got %v expected %v", NUMAPerDev2, expected2)
	}

}
