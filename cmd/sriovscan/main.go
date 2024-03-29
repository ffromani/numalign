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
	"log"
	"os"

	"github.com/ffromani/numalign/internal/pkg/numalign"
	"github.com/ffromani/numalign/internal/pkg/sriovscan"
)

func main() {
	cpusPerNuma, err := numalign.GetCPUsPerNUMANode(numalign.SysDevicesSystemNodeDir)
	if err != nil {
		log.Fatalf("%v", err)
	}

	pciDevs, err := sriovscan.GetPCIDeviceInfo(sriovscan.SysBusPCIDevicesDir)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var sriovDevs []sriovscan.PCIDeviceInfo
	for _, pciDev := range pciDevs {
		if pciDev.IsVFn {
			sriovDevs = append(sriovDevs, pciDev)
		}
	}

	if len(sriovDevs) == 0 {
		os.Exit(0)
	}

	for _, sriovDev := range sriovDevs {
		log.Printf("%s", sriovDev)
	}

	pciPerNuma := sriovscan.CountPCIDevicePerNUMANode(sriovDevs)

	for k := 0; k < len(cpusPerNuma); k++ {
		pciNum := pciPerNuma[k]
		fmt.Printf("%2d: %2d\n", k, pciNum)
	}
}
