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
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/ffromani/cpuset"
)

func printCPUList(cpuList string) {
	cpus, err := cpuset.Parse(cpuList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %q: %v\n", cpuList, err)
		os.Exit(2)
	}

	for _, cpu := range cpus {
		fmt.Printf("%v\n", cpu)
	}

}

func main() {
	var cpuList = flag.StringP("cpu-list", "c", "", "cpulist to split")
	var srcFile = flag.StringP("from-file", "f", "", "read the cpulist to split from the given file")
	flag.Parse()

	if *srcFile == "" && *cpuList == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *cpuList != "" {
		printCPUList(*cpuList)
	}
	if *srcFile != "" {
		var err error
		var data []byte
		if *srcFile == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(*srcFile)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading cpulist from %q: %v\n", *srcFile, err)
			os.Exit(2)
		}
		printCPUList(strings.TrimSpace(string(data)))
	}
}
