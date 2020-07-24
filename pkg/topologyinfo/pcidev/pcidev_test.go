package pcidev

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"testing"

	"github.com/fromanirh/numalign/pkg/fakesysfs"
)

func TestPCIDevsTrivialTree(t *testing.T) {
	base, err := ioutil.TempDir("/tmp", "fakesysfs")
	if err != nil {
		t.Errorf("error creating temp base dir: %v", err)
	}
	fs, err := fakesysfs.NewFakeSysfs(base)
	if err != nil {
		t.Errorf("error creating fakesysfs: %v", err)
	}
	t.Logf("sysfs at %q", fs.Base())

	sysDevs := fs.AddTree("sys", "bus", "pci", "devices")

	attrs := map[string]string{
		"numa_node": "0",
		"class":     "0x020000", // MUST be a network device
		"vendor":    "0x10ec",
		"device":    "0x8168",
	}
	sysDevs.Add("0000:07:00.0", fakesysfs.MakeAttrs(attrs))

	err = fs.Setup()
	if err != nil {
		t.Errorf("error setting up fakesysfs: %v", err)
	}

	pciDevs, err := NewPCIDevices(filepath.Join(fs.Base(), "sys"))
	if err != nil {
		t.Errorf("error in NewPCIDevices: %v", err)
	}

	if len(pciDevs.Items) != 1 {
		t.Errorf("found unexpected amount of PCI devices: %d", len(pciDevs.Items))
	}
	pciDev := pciDevs.Items[0]
	if pciDev.DevClass() != DevClassNetwork {
		t.Errorf("device class mismatch found %x expected %x", pciDev.DevClass(), DevClassNetwork)
	}
	if pciDev.String() != "pci@0000:07:00.0 10ec:8168 numa_node=0 physfn=false vfn=false" {
		t.Errorf("device misdetected: %v", pciDev.String())
	}

	if _, ok := os.LookupEnv("TOPOLOGYINFO_TEST_KEEP_TREE"); ok {
		t.Logf("found environment variable, keeping fake tree")
	} else {
		err = fs.Teardown()
		if err != nil {
			t.Errorf("error tearing down fakesysfs: %v", err)
		}
	}
}
