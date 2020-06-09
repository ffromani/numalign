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

package cmd

import (
	"fmt"

	"github.com/disiqueira/gotree"
	"github.com/spf13/cobra"

	"github.com/fromanirh/numalign/pkg/topologyinfo/pcidev"
)

func showPCIDevs(cmd *cobra.Command, args []string) error {
	pciDevs, err := pcidev.NewPCIDevices("/sys")
	if err != nil {
		return err
	}

	sys := gotree.New(".")
	for nodeId, devInfos := range pciDevs.NUMAPCIDevices {
		numaNode := sys.Add(fmt.Sprintf("numa%02d", nodeId))
		for _, devInfo := range devInfos {
			if sriovInfo, ok := devInfo.(pcidev.SRIOVDeviceInfo); ok && sriovInfo.IsSRIOV() {
				dev := numaNode.Add(fmt.Sprintf("%s %x:%x", sriovInfo.Address(), sriovInfo.Vendor(), sriovInfo.Device()))
				dev.Add(fmt.Sprintf("physfn=%v", sriovInfo.IsPhysFn))
				dev.Add(fmt.Sprintf("vfn=%v", sriovInfo.IsVFn))
			} else {
				numaNode.Add(fmt.Sprintf("%s %x:%x (%x)", devInfo.Address(), devInfo.Vendor(), devInfo.Device(), devInfo.DevClass()))
			}
		}
	}
	fmt.Println(sys.Print())
	return nil
}

func NewPCIDevsCommand() *cobra.Command {
	show := &cobra.Command{
		Use:   "pcidevs",
		Short: "show PCI devices in the system",
		RunE:  showPCIDevs,
		Args:  cobra.NoArgs,
	}
	return show
}
