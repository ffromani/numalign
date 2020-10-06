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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] cpulist\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	var procfsRoot = flag.StringP("procfs", "P", "/proc", "procfs mount point to use.")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	isolCpus, err := cpuset.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %q: %v", args[0], err)
		os.Exit(1)
	}

	irqRoot := filepath.Join(*procfsRoot, "irq")
	var irqViolations []int

	err = filepath.Walk(irqRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		_, file := filepath.Split(path)
		irq, convErr := strconv.Atoi(file)
		if convErr != nil {
			return nil // just skip not-irq-looking dirs
		}

		irqCpuList, readErr := ioutil.ReadFile(filepath.Join(path, "smp_affinity_list"))
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "error reading smp_affinity_list for IRQ %d: %v\n", irq, readErr)
			return nil // keep running
		}

		irqCpus, parseErr := cpuset.Parse(strings.TrimSpace(string(irqCpuList)))

		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "error parsing smp_affinity_list for IRQ %d: %v\n", irq, parseErr)
			return nil // keep running
		}

		cpus := irqCpus.Intersection(isolCpus)
		if cpus.Size() != 0 {
			fmt.Printf("IRQ %3d: can run on %v\n", irq, cpus.ToSlice())
			irqViolations = append(irqViolations, irq)
		}

		return nil
	})
}
