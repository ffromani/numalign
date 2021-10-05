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
 * Copyright 2021 Red Hat, Inc.
 */

package cpusetinfo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const (
	DefaultCGroupsMountPoint string = "/sys/fs/cgroup"
	DefaultProcMountPoint    string = "/proc"
	DefaultSysMountPoint     string = "/sys"
)

const (
	onlineCPUsPath            string = "devices/system/cpu/online"
	cpusetFile                string = "cpuset.cpus"
	threadSiblingListTmplPath string = "devices/system/cpu/cpu%d/topology/thread_siblings_list"
)

const (
	PIDSelf int = 0
)

const (
	cgroupV1 string = "v1"
)

// GetCPUSetForPID retrieves the cpuset allowed for a process, given its pid
func GetCPUSetForPID(fsh FSHandle, pid int) (cpuset.CPUSet, error) {
	cgroupsFile, err := os.Open(cGroupsFileForPID(fsh, pid))
	if err != nil {
		return cpuset.CPUSet{}, err
	}
	defer cgroupsFile.Close()

	cpusPath := ""
	subPath, version := GetCPUSetCGroupPathFromReader(cgroupsFile)
	if subPath == "" {
		cpusPath = filepath.Join(fsh.GetSysMountPoint(), onlineCPUsPath)
	} else {
		switch version {
		case cgroupV1:
			cpusPath = filepath.Join(fsh.GetCGroupsMountPoint(), "cpuset", subPath, cpusetFile)
		default:
			return cpuset.CPUSet{}, fmt.Errorf("detected unsupported cgroup version: %q", version)
		}
	}

	return parseCPUSetFile(cpusPath)
}

type ThreadSiblingMap struct {
	siblings map[int][]int
	fsh      FSHandle
}

func NewThreadSiblingMap(fsh FSHandle) *ThreadSiblingMap {
	return &ThreadSiblingMap{
		siblings: make(map[int][]int),
		fsh:      fsh,
	}
}

func (tsm *ThreadSiblingMap) SetCPUSiblings(cpu int, siblings []int) *ThreadSiblingMap {
	tsm.siblings[cpu] = siblings
	return tsm
}

func (tsm *ThreadSiblingMap) ForCPU(cpu int) ([]int, error) {
	if val, ok := tsm.siblings[cpu]; ok {
		return val, nil
	}

	ts, err := parseCPUSetFile(filepath.Join(tsm.fsh.GetSysMountPoint(), fmt.Sprintf(threadSiblingListTmplPath, cpu)))
	if err != nil {
		return []int{}, err
	}
	val := ts.ToSliceNoSort()
	tsm.siblings[cpu] = val
	return val, nil
}

// CheckCPUSetIsSiblingAligned tells if a given cpuset is composed only by thread siblings sets,
// IOW if core-level noisy neighbours are possible or not. Returns the misaligned CPU IDs.
func (tsm ThreadSiblingMap) CheckCPUSetAligned(cpus cpuset.CPUSet) (cpuset.CPUSet, error) {
	misaligned := cpuset.CPUSet{}
	reconstructed := cpuset.NewBuilder()
	for _, cpuId := range cpus.ToSliceNoSort() {
		ts, err := tsm.ForCPU(cpuId)
		if err != nil {
			return misaligned, err
		}
		reconstructed.Add(ts...)
	}

	found := reconstructed.Result()
	if found.Equals(cpus) {
		return misaligned, nil
	}

	builder := cpuset.NewBuilder()
	for _, extraCpu := range found.Difference(cpus).ToSliceNoSort() {
		cpuIds, err := tsm.ForCPU(extraCpu)
		if err != nil {
			return misaligned, err
		}
		for _, cpuId := range cpuIds {
			if !cpus.Contains(cpuId) {
				continue
			}
			builder.Add(cpuId)
		}
	}
	return builder.Result(), nil
}

func GetCPUSetCGroupPathFromReader(r io.Reader) (string, string) {
	return getCPUSetCGroupPathFromReaderV1(r), cgroupV1
}

func getCPUSetCGroupPathFromReaderV1(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		entry := strings.TrimSpace(scanner.Text())
		if !strings.Contains(entry, "cpuset") {
			continue
		}
		// entry format is "number:name:path"
		items := strings.Split(entry, ":")
		if len(items) != 3 {
			// how come?
			continue
		}
		return items[2]
	}
	return ""
}

type FSHandle struct {
	CGroupsMountPoint string
	ProcMountPoint    string
	SysMountPoint     string
}

func (fsh FSHandle) GetCGroupsMountPoint() string {
	if fsh.CGroupsMountPoint == "" {
		return DefaultCGroupsMountPoint
	}
	return fsh.CGroupsMountPoint
}

func (fsh FSHandle) GetProcMountPoint() string {
	if fsh.ProcMountPoint == "" {
		return DefaultProcMountPoint
	}
	return fsh.ProcMountPoint
}

func (fsh FSHandle) GetSysMountPoint() string {
	if fsh.SysMountPoint == "" {
		return DefaultSysMountPoint
	}
	return fsh.SysMountPoint
}

func cGroupsFileForPID(fsh FSHandle, pid int) string {
	pidStr := "self"
	if pid > 0 && pid != PIDSelf {
		pidStr = fmt.Sprintf("%d", pid)
	}
	return filepath.Join(fsh.GetProcMountPoint(), pidStr, "cgroup")
}

func parseCPUSetFile(cpusPath string) (cpuset.CPUSet, error) {
	data, err := os.ReadFile(cpusPath)
	if err != nil {
		return cpuset.CPUSet{}, err
	}
	return cpuset.Parse(strings.TrimSpace(string(data)))
}
