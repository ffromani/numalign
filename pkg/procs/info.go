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

package procs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ffromani/cpuset"
)

type Info struct {
	Pid      int32
	Affinity []int
}

func All(procfsRoot string) ([]Info, error) {
	infos := []Info{}
	entries, err := ioutil.ReadDir(procfsRoot)
	if err != nil {
		return infos, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(entry.Name()); err != nil {
			// doesn't look like a pid
			continue
		}
		info, err := parseProcStatus(filepath.Join(procfsRoot, entry.Name(), "status"))
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func FromPID(procfsRoot string, pid int32) (Info, error) {
	return parseProcStatus(filepath.Join(procfsRoot, fmt.Sprintf("%d", pid), "status"))
}

func parseProcStatus(path string) (Info, error) {
	info := Info{}
	file, err := os.Open(path)
	if err != nil {
		return info, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Pid:") {
			items := strings.SplitN(line, ":", 2)
			pid, err := strconv.Atoi(strings.TrimSpace(items[1]))
			if err != nil {
				return info, err
			}
			info.Pid = int32(pid)
		}
		if strings.HasPrefix(line, "Cpus_allowed_list:") {
			items := strings.SplitN(line, ":", 2)
			cpuIDs, err := cpuset.Parse(strings.TrimSpace(items[1]))
			if err != nil {
				return info, err
			}
			info.Affinity = cpuIDs
		}
	}

	return info, scanner.Err()
}
