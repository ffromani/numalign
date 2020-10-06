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
	"sort"
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
	var checkEffective = flag.BoolP("effective-affinity", "E", false, "check effective affinity.")
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

	files, err := ioutil.ReadDir(irqRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %q: %v", irqRoot, err)
		os.Exit(1)
	}

	var irqs []int
	for _, file := range files {
		irq, err := strconv.Atoi(file.Name())
		if err != nil {
			continue // just skip not-irq-looking dirs
		}
		irqs = append(irqs, irq)
	}

	sort.Ints(irqs)

	affinityListFile := "smp_affinity_list"
	if *checkEffective {
		affinityListFile = "effective_affinity_list"
	}

	for _, irq := range irqs {
		irqDir := filepath.Join(irqRoot, fmt.Sprintf("%d", irq))

		irqCpuList, err := ioutil.ReadFile(filepath.Join(irqDir, affinityListFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %q for IRQ %d: %v\n", affinityListFile, irq, err)
			continue // keep running
		}

		irqCpus, err := cpuset.Parse(strings.TrimSpace(string(irqCpuList)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing %q for IRQ %d: %v\n", affinityListFile, irq, err)
			continue // keep running
		}

		cpus := irqCpus.Intersection(isolCpus)
		if cpus.Size() != 0 {
			source := findSourceForIRQ(*procfsRoot, irq)
			fmt.Printf("IRQ %3d [%24s]: can run on %v\n", irq, source, cpus.ToSlice())
			irqViolations = append(irqViolations, irq)
		}
	}

	if len(irqViolations) > 0 {
		os.Exit(1)
	}
}

func findSourceForIRQ(procfs string, irq int) string {
	irqDir := filepath.Join(procfs, "irq", fmt.Sprintf("%d", irq))
	files, err := ioutil.ReadDir(irqDir)
	if err != nil {
		return "MISSING"
	}
	for _, file := range files {
		if file.IsDir() {
			return file.Name()
		}
	}
	return ""
}
