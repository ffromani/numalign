package fakesysfs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"testing"
)

type devData struct {
	addr     string
	numaNode string
}

func TestSinglePCIDevice(t *testing.T) {
	base, err := ioutil.TempDir("/tmp", "fakesysfs")
	if err != nil {
		t.Errorf("error creating temp base dir: %v", err)
	}
	fs, err := NewFakeSysfs(base)
	if err != nil {
		t.Errorf("error creating fakesysfs: %v", err)
	}
	root := fs.Root()
	if root == nil {
		t.Errorf("nil fakesysfs root")
	}
	t.Logf("sysfs at %q", fs.Base())

	fakeDevs := []devData{
		{
			addr:     "0000:00:01.0",
			numaNode: "0",
		},
		{
			addr:     "0000:00:02.0",
			numaNode: "1",
		},
	}

	devs := fs.AddTree("sys", "bus", "pci", "devices")
	for _, fakeDev := range fakeDevs {
		devs.Add(fakeDev.addr, map[string]string{
			"numa_node": fmt.Sprintf("%s\n", fakeDev.numaNode),
		})
	}
	err = fs.Setup()
	if err != nil {
		t.Errorf("error setting up fakesysfs: %v", err)
	}

	for _, fakeDev := range fakeDevs {
		data, err := ioutil.ReadFile(filepath.Join(base, "sys", "bus", "pci", "devices", fakeDev.addr, "numa_node"))
		if err != nil {
			t.Errorf("error reading back %q: %v", fakeDev.addr, err)
		}
		val := strings.TrimSpace(string(data))
		if val != fakeDev.numaNode {
			t.Errorf("value mismatch for %q expected %q got %q", fakeDev.addr, fakeDev.numaNode, val)
		}
	}

	err = fs.Teardown()
	if err != nil {
		t.Errorf("error tearing down fakesysfs: %v", err)
	}
}
