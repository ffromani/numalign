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

package numa

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	fakesysfs "github.com/ffromani/numalign/pkg/topologyinfo/sysfs/fake"
)

func TestReadNUMAINfo(t *testing.T) {

	expected := Nodes{
		Online:           []int{0, 2, 3},
		Possible:         []int{0, 1, 2, 3},
		WithCPU:          []int{0, 1},
		WithMemory:       []int{2, 3},
		WithNormalMemory: []int{0, 1},
	}

	data := map[string]string{
		"online":            "0,2,3",
		"possible":          "0,1,2,3",
		"has_cpu":           "0,1",
		"has_memory":        "2,3",
		"has_normal_memory": "0,1",
	}

	base, err := ioutil.TempDir("/tmp", "fakesysfs")
	if err != nil {
		t.Errorf("error creating temp base dir: %v", err)
	}
	fs, err := fakesysfs.NewFakeSysfs(base)
	if err != nil {
		t.Errorf("error creating fakesysfs: %v", err)
	}
	t.Logf("sysfs at %q", fs.Base())

	sysDevs := fs.AddTree("sys", "devices")
	devSys := sysDevs.Add("system", nil)
	devSys.Add("node", data)

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

	info, err := NewNodesFromSysFS(filepath.Join(fs.Base(), "sys"))
	if err != nil {
		t.Errorf("error in NewNodesFromSysFS: %v", err)
	}

	if !reflect.DeepEqual(info, expected) {
		t.Errorf("data mismatch found %v expected %v", info, expected)
	}
}
