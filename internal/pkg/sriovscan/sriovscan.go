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
	NUMANode int
	IsPhysFn bool
	IsVFn    bool
}

func (pdi PCIDeviceInfo) String() string {
	return fmt.Sprintf("pci@%s numa_node=%d physfn=%v vfn=%v", pdi.Address, pdi.NUMANode, pdi.IsPhysFn, pdi.IsVFn)
}

func GetPCIDeviceInfo(sysPCIDir string) ([]PCIDeviceInfo, error) {
	var pciDevs []PCIDeviceInfo

	entries, err := ioutil.ReadDir(SysBusPCIDevicesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		isPhysFn := false
		isVFn := false
		if _, err := os.Stat(filepath.Join(sysPCIDir, entry.Name(), "sriov_numvfs")); err == nil {
			isPhysFn = true
		} else if !os.IsNotExist(err) {
			// unexpected error. Bail out
			return nil, err
		}
		if _, err := os.Stat(filepath.Join(sysPCIDir, entry.Name(), "physfn")); err == nil {
			isVFn = true
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

		pciDevs = append(pciDevs, PCIDeviceInfo{
			Address:  entry.Name(),
			NUMANode: nodeNum,
			IsPhysFn: isPhysFn,
			IsVFn:    isVFn,
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
