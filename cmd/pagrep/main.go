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

	"github.com/ffromani/numalign/pkg/procs"
	k8scpuset "k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [cpulist]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	var procfsRoot = flag.StringP("procfs", "P", "/proc", "procfs mount point to use.")
	flag.Parse()

	isolCpuList := "0-65535" // "everything"
	args := flag.Args()
	if len(args) == 1 {
		isolCpuList = args[0]
	}

	isolCpus, err := k8scpuset.Parse(isolCpuList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %q: %v", isolCpuList, err)
		os.Exit(1)
	}

	procInfos, err := procs.All(*procfsRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting process infos from %q: %v", *procfsRoot, err)
		os.Exit(1)
	}

	for _, procInfo := range procInfos {
		procCpus := k8scpuset.NewCPUSet(procInfo.Affinity...)

		cpus := procCpus.Intersection(isolCpus)
		if cpus.Size() != 0 {
			fmt.Printf("PID %6d can run on %v\n", procInfo.Pid, cpus.ToSlice())
		}
	}
}
