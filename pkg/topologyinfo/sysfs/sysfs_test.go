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

package sysfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fakesysfs "github.com/fromanirh/numalign/pkg/topologyinfo/sysfs/fake"
)

func TestReadSingleNuma(t *testing.T) {
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
	cpuList := "0-7\n"

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

	sysRoot := filepath.Join(fs.Base(), "sys")
	foundCpus, err := New(sysRoot).ForNode(0).ReadList("cpulist")
	if err != nil {
		t.Errorf("unexpected error reading cpulist for node 0: %v", err)
	}
	if !reflect.DeepEqual(foundCpus, allCpus) {
		t.Errorf("found cpus %v expected %v", foundCpus, allCpus)
	}

	foundPhysPackageID, err := New(sysRoot).ForCPU(0).Join("topology").ReadFile("physical_package_id")
	if err != nil {
		t.Errorf("unexpected error reading physical package id for cpu 0: %v", err)
	}
	val := strings.TrimSpace(foundPhysPackageID)
	if val != "0" {
		t.Errorf("found physical package id %v expected %v", val, "0")
	}
}

func TestReadWrongData(t *testing.T) {
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
	cpuList := "0-7\n"

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

	sysRoot := filepath.Join(fs.Base(), "sys")
	foundCpus, err := New(sysRoot).ForNode(0).ReadList("cpulist")
	if err != nil {
		t.Errorf("unexpected error reading cpulist for node 0: %v", err)
	}
	if !reflect.DeepEqual(foundCpus, allCpus) {
		t.Errorf("found cpus %v expected %v", foundCpus, allCpus)
	}

	_, err = New(sysRoot).ForCPU(0).Join("topology").ReadFile("does_not_exist")
	if err == nil {
		t.Errorf("missing expected error reading unexistent file")
	}

	_, err = New(sysRoot).ForCPU(0).Join("topology").ReadList("does_not_exist")
	if err == nil {
		t.Errorf("missing expected error reading unexistent file as list")
	}
}
