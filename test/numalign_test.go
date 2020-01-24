package cmd_test

import (
	"testing"

	"github.com/fromanirh/numalign/cmd"
)

func TestResourceString(t *testing.T) {
	var r1 cmd.Resources
	s1 := r1.String()
	if s1 != "" {
		t.Errorf("empty Resource has unexpected output: %s", s1)
	}

	r2 := cmd.Resources{
		CPUToNUMANode: map[int]int{
			0: 0,
			1: 0,
		},
		PCIDevsToNUMANode: map[string]int{
			"3c:00.0": 0,
		},
	}
	s2 := r2.String()
	expected := `CPU cpu#000=00
CPU cpu#001=00
PCI 3c:00.0=00
`
	if s2 != expected {
		t.Errorf("initialzed Resource has unexpected output: %s (expected %s)", s2, expected)
	}

}
