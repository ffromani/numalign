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

package cpus

import (
	"fmt"
	"io"
	"strings"

	"github.com/fromanirh/cpuset"
)

func summarizeCPUIdList(data map[int]CPUIdList) string {
	ref := 0
	var items []string
	for cpuID, cpuList := range data {
		cur := len(cpuList)
		if ref == 0 {
			ref = cur
		} else if ref != cur {
			items = append(items, fmt.Sprintf("core%d=%d", cpuID, cur))
		}
	}
	if len(items) > 0 {
		return strings.Join(items, ",")
	}
	return fmt.Sprintf("%d", ref)
}

func MakeSummary(cpuInfos *CPUs, w io.Writer) {
	fmt.Fprintf(w, "CPU(s):\t%d\n", len(cpuInfos.Present))
	fmt.Fprintf(w, "Present CPU(s) list:\t%s\n", cpuset.Unparse(cpuInfos.Present))
	fmt.Fprintf(w, "On-line CPU(s) list:\t%s\n", cpuset.Unparse(cpuInfos.Online))
	fmt.Fprintf(w, "Thread(s) per core:\t%s\n", summarizeCPUIdList(cpuInfos.CoreCPUs))
	fmt.Fprintf(w, "Core(s) per socket:\t%s\n", summarizeCPUIdList(cpuInfos.PackageCPUs))
	fmt.Fprintf(w, "Socket(s):\t%d\n", cpuInfos.Packages)
	fmt.Fprintf(w, "NUMA node(s):\t%d\n", len(cpuInfos.NUMANodes))
	for _, idx := range cpuInfos.NUMANodes {
		fmt.Fprintf(w, "NUMA node%d CPU(s):\t%s\n", idx, cpuset.Unparse(cpuInfos.NUMANodeCPUs[idx]))
	}
}
