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
	"strings"
	"testing"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const cgroupData string = `12:cpuset:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
11:freezer:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
10:memory:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
9:misc:/
8:blkio:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
7:perf_event:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
6:pids:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
5:devices:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
4:cpu,cpuacct:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
3:hugetlb:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
2:net_cls,net_prio:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
1:name=systemd:/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854
0::/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/system.slice/containerd.service`

func NewTestTopology() *ThreadSiblingMap {
	tsm := NewThreadSiblingMap(FSHandle{})
	tsm.SetCPUSiblings(0, []int{0, 16})
	tsm.SetCPUSiblings(1, []int{1, 17})
	tsm.SetCPUSiblings(2, []int{2, 18})
	tsm.SetCPUSiblings(3, []int{3, 19})
	tsm.SetCPUSiblings(4, []int{4, 20})
	tsm.SetCPUSiblings(5, []int{5, 21})
	tsm.SetCPUSiblings(6, []int{6, 22})
	tsm.SetCPUSiblings(7, []int{7, 23})
	tsm.SetCPUSiblings(8, []int{8, 24})
	tsm.SetCPUSiblings(9, []int{9, 25})
	tsm.SetCPUSiblings(10, []int{10, 26})
	tsm.SetCPUSiblings(11, []int{11, 27})
	tsm.SetCPUSiblings(12, []int{12, 28})
	tsm.SetCPUSiblings(13, []int{13, 29})
	tsm.SetCPUSiblings(14, []int{14, 30})
	tsm.SetCPUSiblings(15, []int{15, 31})
	tsm.SetCPUSiblings(16, []int{0, 16})
	tsm.SetCPUSiblings(17, []int{1, 17})
	tsm.SetCPUSiblings(18, []int{2, 18})
	tsm.SetCPUSiblings(19, []int{3, 19})
	tsm.SetCPUSiblings(20, []int{4, 20})
	tsm.SetCPUSiblings(21, []int{5, 21})
	tsm.SetCPUSiblings(22, []int{6, 22})
	tsm.SetCPUSiblings(23, []int{7, 23})
	tsm.SetCPUSiblings(24, []int{8, 24})
	tsm.SetCPUSiblings(25, []int{9, 25})
	tsm.SetCPUSiblings(26, []int{10, 26})
	tsm.SetCPUSiblings(27, []int{11, 27})
	tsm.SetCPUSiblings(28, []int{12, 28})
	tsm.SetCPUSiblings(29, []int{13, 29})
	tsm.SetCPUSiblings(30, []int{14, 30})
	tsm.SetCPUSiblings(31, []int{15, 31})
	return tsm
}

func TestGetCPUSetCGroupPathFromReader(t *testing.T) {
	testCases := []struct {
		description     string
		data            string
		expectedVersion string
		expectedPath    string
	}{
		{
			description:     "missing v1 data",
			data:            "",
			expectedVersion: cgroupV1,
			expectedPath:    "",
		},

		{
			description:     "valid v1 data",
			data:            cgroupData,
			expectedVersion: cgroupV1,
			expectedPath:    "/docker/95b99ca10ff72f086a51561b32957244ef498e88d5564a11fdbae039cc42d581/kubelet/kubepods/podb1c81bdc-1bc5-4d39-a173-b74598538a91/741e4d6c8494d2492df382a0c3f765c424bd784869fdf5a399cfbeba71e11854",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			subPath, version := GetCPUSetCGroupPathFromReader(strings.NewReader(testCase.data))
			if version != testCase.expectedVersion {
				t.Errorf("detected unexpected cgroup version %q", version)
			}
			if subPath != testCase.expectedPath {
				t.Errorf("detected unexpected cgroup subPath %q", subPath)
			}

		})
	}
}

func TestCheckCPUSetAligned(t *testing.T) {
	tsm := NewTestTopology()
	testCases := []struct {
		description        string
		cpus               cpuset.CPUSet
		expectedMisaligned cpuset.CPUSet
		expectedError      bool
	}{
		{
			description:        "single cpu",
			cpus:               cpuset.NewCPUSet(4),
			expectedMisaligned: cpuset.NewCPUSet(4),
		},
		{
			description:        "requesting even number of cpus",
			cpus:               cpuset.NewCPUSet(2, 3, 18, 19),
			expectedMisaligned: cpuset.CPUSet{},
		},
		{
			description:        "requesting odd number of cpus",
			cpus:               cpuset.NewCPUSet(2, 3, 18),
			expectedMisaligned: cpuset.NewCPUSet(3),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			misaligned, err := tsm.CheckCPUSetAligned(testCase.cpus)
			if !misaligned.Equals(testCase.expectedMisaligned) {
				t.Errorf("unexpected misaligned cpus: got %v expected %v", misaligned.ToSlice(), testCase.expectedMisaligned.ToSlice())
			}
			gotError := err != nil
			if gotError != testCase.expectedError {
				t.Errorf("unexpected error: got %v expected %v", err, testCase.expectedError)
			}
		})
	}
}
