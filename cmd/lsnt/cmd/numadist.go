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
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/ffromani/numalign/pkg/topologyinfo/numa"
	"github.com/ffromani/numalign/pkg/topologyinfo/numa/distances"
)

func nodeIDs(nodeIDs []int) []string {
	ret := []string{}
	for _, nodeID := range nodeIDs {
		ret = append(ret, fmt.Sprintf("%3d", nodeID))
	}
	return ret
}

func showNUMADist(cmd *cobra.Command, args []string) error {
	nodes, err := numa.NewNodesFromSysFS(opts.sysFSRoot)
	if err != nil {
		return err
	}
	dists, err := distances.NewDistancesFromSysfs(opts.sysFSRoot)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 4, 1, ' ', tabwriter.AlignRight)

	// hack to fix the alignment
	fmt.Fprintf(w, "   %s\n", strings.Join(nodeIDs(nodes.Online), "\t")) // header
	for _, nodeIDFrom := range nodes.Online {
		fmt.Fprintf(w, "%d:", nodeIDFrom)
		for _, nodeIDTo := range nodes.Online {
			val, err := dists.BetweenNodes(nodeIDFrom, nodeIDTo)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "\t%3d", val)
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return nil
}

func newNUMADistCommand() *cobra.Command {
	show := &cobra.Command{
		Use:   "numadist",
		Short: "show NUMA distances",
		RunE:  showNUMADist,
		Args:  cobra.NoArgs,
	}
	return show
}
