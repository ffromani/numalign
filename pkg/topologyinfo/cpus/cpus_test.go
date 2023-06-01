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
	"io/ioutil"
	"os"
	"path/filepath"

	"testing"

	"github.com/ffromani/cpuset"
	fakesysfs "github.com/ffromani/numalign/pkg/topologyinfo/sysfs/fake"
	"github.com/google/go-cmp/cmp"
)

func TestCPUsSingleNuma(t *testing.T) {
	base, err := ioutil.TempDir("/tmp", "fakesysfs")
	if err != nil {
		t.Errorf("error creating temp base dir: %v", err)
	}
	fs, err := fakesysfs.NewFakeSysfs(base)
	if err != nil {
		t.Errorf("error creating fakesysfs: %v", err)
	}
	t.Logf("sysfs at %q", fs.Base())

	allCpus := []int{0, 1, 2, 3, 4, 5, 6, 7}
	cpuList := cpuset.Unparse(allCpus) + "\n"

	sysDevs := fs.AddTree("sys", "devices")
	devSys := sysDevs.Add("system", nil)
	devNode := devSys.Add("node", map[string]string{
		"online": "0",
	})
	devNode.Add("node0", map[string]string{
		"cpulist": cpuList,
	})
	devCpu := devSys.Add("cpu", map[string]string{
		"present": cpuList,
		"online":  cpuList,
	})
	for _, cpuID := range allCpus {
		devCpu.Add(fmt.Sprintf("cpu%d", cpuID), nil).Add("topology", map[string]string{
			"thread_siblings_list": cpuList,
			"core_siblings_list":   cpuList,
			"physical_package_id":  "0\n",
		})
	}

	err = fs.Setup()
	if err != nil {
		t.Errorf("error setting up fakesysfs: %v", err)
	}
	defer func() {
		if _, ok := os.LookupEnv("TOPOLOGYINFO_TEST_KEEP_TREE"); ok {
			t.Logf("found environment variable, keeping fake tree")
		} else {
			err = fs.Teardown()
			if err != nil {
				t.Errorf("error tearing down fakesysfs: %v", err)
			}
		}
	}()

	cpus, err := NewCPUs(filepath.Join(fs.Base(), "sys"))
	if err != nil {
		t.Errorf("error in NewCPU: %v", err)
	}
	if len(cpus.NUMANodes) != 1 || len(cpus.NUMANodeCPUs) != 1 {
		t.Errorf("NUMA Nodes miscount: expected %d detected %d/%d", 1, len(cpus.NUMANodes), len(cpus.NUMANodeCPUs))
	}

	testingCpus := CPUIdList(allCpus)
	if !cmp.Equal(cpus.Present, testingCpus) {
		t.Errorf("not all cpus present %v vs %v", cpus.Present, testingCpus)
	}
	if !cmp.Equal(cpus.Online, testingCpus) {
		t.Errorf("not all cpus online: %v vs %v", cpus.Online, testingCpus)
	}
	if !cmp.Equal(cpus.NUMANodeCPUs[0], testingCpus) {
		t.Errorf("not all cpus on NUMA#0: %v vs %v", cpus.NUMANodeCPUs[0], testingCpus)
	}
}
