package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fromanirh/numalign/internal/pkg/numalign"
	"github.com/fromanirh/numalign/internal/pkg/sriovscan"
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
		if pciDev.IsSRIOV {
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
