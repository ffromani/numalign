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
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/ffromani/numalign/pkg/cpusetinfo"
)

type result struct {
	Aligned        bool  `json:"aligned"`
	CPUsAllowed    []int `json:"cpus_allowed"`
	CPUsMisaligned []int `json:"cpus_misaligned"`
	Pid            int   `json:"pid"`
}

func main() {
	flag.Parse()
	pids := flag.Args()

	if len(pids) != 0 && len(pids) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	pid := 0
	if len(pids) == 1 && pids[0] != "self" {
		v, err := strconv.Atoi(pids[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad argument %q: %v\n", pids[0], err)
			os.Exit(2)
		}
		pid = v
	}

	fsh := cpusetinfo.FSHandle{}

	cpus, err := cpusetinfo.GetCPUSetForPID(fsh, pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot fetch the cpuset for pid %d: %v\n", pid, err)
		os.Exit(2)
	}

	tsm := cpusetinfo.NewThreadSiblingMap(fsh)

	misaligned, err := tsm.CheckCPUSetAligned(cpus)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot check the cpuset for pid %d: %v\n", pid, err)
		os.Exit(4)
	}

	err = json.NewEncoder(os.Stdout).Encode(result{
		Aligned:        misaligned.Size() == 0,
		CPUsAllowed:    cpus.ToSlice(),
		CPUsMisaligned: misaligned.ToSlice(),
		Pid:            pid,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot encode the result: %v\n", err)
		os.Exit(8)
	}
}
