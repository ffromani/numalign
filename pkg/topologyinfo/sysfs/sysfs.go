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
	"path/filepath"
	"strings"

	"github.com/fromanirh/cpuset"
)

/*
 * keep this handy:
 * https://www.kernel.org/doc/html/latest/admin-guide/cputopology.html
 */
const (
	PathDevsSysCPU  = "devices/system/cpu"
	PathDevsSysNode = "devices/system/node"
)

type Path struct {
	base string
}

func New(base string) Path {
	return Path{
		base: base,
	}
}

func (p Path) Join(extra string) Path {
	return Path{
		base: filepath.Join(p.base, extra),
	}
}

func (p Path) ReadFile(name string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(p.base, name))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (p Path) ReadList(name string) ([]int, error) {
	data, err := p.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return cpuset.Parse(data)
}

func (p Path) ForNode(nodeID int) Path {
	return Path{
		base: filepath.Join(p.base, PathDevsSysNode, fmt.Sprintf("node%d", nodeID)),
	}
}

func (p Path) ForCPU(cpuID int) Path {
	return Path{
		base: filepath.Join(p.base, PathDevsSysCPU, fmt.Sprintf("cpu%d", cpuID)),
	}
}
