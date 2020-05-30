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
	"text/tabwriter"

	"github.com/fromanirh/cpuset"
	"github.com/fromanirh/numalign/pkg/k8sresource/cpus"
)

func ExpectNoError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func main() {
	cpuRes, err := cpus.NewCPUs("/sys")
	ExpectNoError(err)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintf(w, "CPU(s):\t%d\n", cpuRes.Total)
	fmt.Fprintf(w, "On-line CPU(s) list:\t%s\n", cpuset.Unparse(cpuRes.Online))
	fmt.Fprintf(w, "Thread(s) per core:\t%d\n", cpuRes.ThreadsPerCore)
	fmt.Fprintf(w, "Core(s) per socket:\t%d\n", cpuRes.CoresPerSocket)
	fmt.Fprintf(w, "Socket(s):\t%d\n", cpuRes.Sockets)
	fmt.Fprintf(w, "NUMA node(s):\t%d\n", cpuRes.NUMANodes)
	for i := 0; i < cpuRes.NUMANodes; i++ {
		fmt.Fprintf(w, "NUMA node%d CPU(s):\t%s\n", cpuset.Unparse(cpuRes.NUMANodeCPUs[i]))
	}
	w.Flush()
}
