package fake

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"testing"
)

type devData struct {
	addr  string
	attrs map[string]string
}

func TestPCIDevices(t *testing.T) {
	DebugLog = log.Printf

	base, err := ioutil.TempDir("/tmp", "fakesysfs")
	if err != nil {
		t.Errorf("error creating temp base dir: %v", err)
	}
	fs, err := NewFakeSysfs(base)
	if err != nil {
		t.Errorf("error creating fakesysfs: %v", err)
	}
	t.Logf("sysfs at %q", fs.Base())

	fakeDevs := []devData{
		{
			addr: "0000:00:01.0",
			attrs: map[string]string{
				"numa_node":     "0\n",
				"vendor":        "0x8086\n",
				"device":        "0xcafe\n",
				"local_cpulist": "0,42\n",
			},
		},
		{
			addr: "0000:00:02.0",
			attrs: map[string]string{
				"numa_node": "1\n",
			},
		},
	}

	devs := fs.AddTree("sys", "bus", "pci", "devices")
	for _, fakeDev := range fakeDevs {
		devs.Add(fakeDev.addr, fakeDev.attrs)
	}
	err = fs.Setup()
	if err != nil {
		t.Errorf("error setting up fakesysfs: %v", err)
	}

	defer func() {
		if _, ok := os.LookupEnv("FAKESYS_TEST_KEEP_TREE"); ok {
			t.Logf("found environment variable, keeping fake tree")
		} else {
			err = fs.Teardown()
			if err != nil {
				t.Errorf("error tearing down fakesysfs: %v", err)
			}
		}
	}()

	checkPath(t, fakeDevs, filepath.Join(fs.Base(), "sys", "bus", "pci", "devices"))
}

func checkPath(t *testing.T, fakeDevs []devData, subPath string) {
	for _, fakeDev := range fakeDevs {
		devPath := filepath.Join(subPath, fakeDev.addr)

		attrs, err := readNodeAttrs(devPath)
		if err != nil {
			t.Errorf("error reading back %q: %v", fakeDev.addr, err)
		}

		for key, value := range fakeDev.attrs {
			val, ok := attrs[key]
			if !ok || val != value {
				t.Errorf("value mismatch for %q key %q expected %q got %q", fakeDev.addr, key, value, val)
			}
		}
	}
}

func readNodeAttrs(leafPath string) (map[string]string, error) {
	files, err := ioutil.ReadDir(leafPath)
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]string)
	for _, f := range files {
		attrPath := filepath.Join(leafPath, f.Name())
		data, err := ioutil.ReadFile(attrPath)
		if err != nil {
			return attrs, err
		}

		attrs[f.Name()] = string(data)
	}
	return attrs, nil
}
