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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/ffromani/numalign/pkg/topologyinfo/pcidev"
)

type sriovInfo struct {
	PFDev  pcidev.SRIOVDeviceInfo
	VFDevs []pcidev.SRIOVDeviceInfo
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] physfn_pci_addr\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	var sysfsRoot = flag.StringP("sysfs", "S", "/sys", "sysfs mount point to use.")
	var numaNode = flag.IntP("numa-node", "N", 0, "numa node to pin to")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	pfPCIAddr := args[0]

	devInfos, err := pcidev.NewPCIDevices(*sysfsRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PCI device listing failed: %v", err)
		os.Exit(1)
	}

	var info sriovInfo
	havePF := false
	for _, devInfo := range devInfos.Items {
		if devInfo.DevClass() != pcidev.DevClassNetwork {
			continue
		}
		sriovDev, ok := devInfo.(pcidev.SRIOVDeviceInfo)
		if !ok {
			continue
		}

		if sriovDev.IsPhysFn && (sriovDev.Address() == pfPCIAddr || sriovDev.DevAddress() == pfPCIAddr) {
			havePF = true
			info.PFDev = sriovDev
		}
	}

	if !havePF {
		fmt.Fprintf(os.Stderr, "phsyfn %q not found in the system\n", pfPCIAddr)
		os.Exit(0)
	}

	for _, devInfo := range devInfos.Items {
		if devInfo.DevClass() != pcidev.DevClassNetwork {
			continue
		}
		sriovDev, ok := devInfo.(pcidev.SRIOVDeviceInfo)
		if !ok {
			continue
		}

		if sriovDev.IsVFn && (sriovDev.ParentFn == info.PFDev.Address()) {
			info.VFDevs = append(info.VFDevs, sriovDev)
		}
	}

	if info.PFDev.NUMANode() != *numaNode {
		fmt.Printf("echo %d > %s\n", *numaNode, filepath.Join(info.PFDev.SysfsPath(), "numa_node"))
	}
	for _, vfDev := range info.VFDevs {
		if vfDev.NUMANode() != *numaNode {
			fmt.Printf("echo %d > %s\n", *numaNode, filepath.Join(vfDev.SysfsPath(), "numa_node"))
		}
	}
}
