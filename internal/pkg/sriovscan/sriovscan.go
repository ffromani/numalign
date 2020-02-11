package sriovscan

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	SysBusPCIDevicesDir = "/sys/bus/pci/devices/"
)

type PCIDeviceInfo struct {
	Address  string
	IsSRIOV  bool
	NUMANode int
}

func (pdi PCIDeviceInfo) String() string {
	return fmt.Sprintf("pci@%s numa_node=%d sriov=%v", pdi.Address, pdi.NUMANode, pdi.IsSRIOV)
}

func GetPCIDeviceInfo(sysPCIDir string) ([]PCIDeviceInfo, error) {
	var pciDevs []PCIDeviceInfo

	entries, err := ioutil.ReadDir(SysBusPCIDevicesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		isSRIOV := false
		if _, err := os.Stat(filepath.Join(sysPCIDir, entry.Name(), "sriov_numvfs")); err == nil {
			isSRIOV = true
		} else if !os.IsNotExist(err) {
			// unexpected error. Bail out
			return nil, err
		}

		content, err := ioutil.ReadFile(filepath.Join(sysPCIDir, entry.Name(), "numa_node"))
		if err != nil {
			return nil, err
		}
		nodeNum, err := strconv.Atoi(strings.TrimSpace(string(content)))
		if err != nil {
			return nil, err
		}

		// XXX
		if nodeNum == -1 {
			nodeNum = 0
		}

		pciDevs = append(pciDevs, PCIDeviceInfo{
			Address:  entry.Name(),
			IsSRIOV:  isSRIOV,
			NUMANode: nodeNum,
		})
	}

	return pciDevs, nil
}

func CountPCIDevicePerNUMANode(pciDevs []PCIDeviceInfo) map[int]int {
	pciPerNuma := make(map[int]int)
	for _, pciDev := range pciDevs {
		pciPerNuma[pciDev.NUMANode] += 1
	}
	return pciPerNuma
}
