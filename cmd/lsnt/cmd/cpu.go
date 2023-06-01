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
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/ffromani/numalign/pkg/topologyinfo/cpus"
)

func showCPU(cmd *cobra.Command, args []string) error {
	cpuInfos, err := cpus.NewCPUs(opts.sysFSRoot)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	cpus.MakeSummary(cpuInfos, w)
	w.Flush()
	return nil
}

func newCPUCommand() *cobra.Command {
	show := &cobra.Command{
		Use:   "cpu",
		Short: "show cpu details like lscpu(1)",
		RunE:  showCPU,
		Args:  cobra.NoArgs,
	}
	return show
}
