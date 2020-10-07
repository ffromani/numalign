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

package softirqs

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// presented in kernel order
func Names() []string {
	return []string{"HI", "TIMER", "NET_TX", "NET_RX", "BLOCK", "IRQ_POLL", "TASKLET", "SCHED", "HRTIMER", "RCU"}
}

type Info struct {
	// Count online cpus
	CPUs int
	// softirq -> percpu count
	Counters map[string][]uint64
}

func ReadInfo(rd io.Reader) (*Info, error) {
	src := bufio.NewScanner(rd)
	src.Scan()
	cpus := strings.Fields(src.Text())
	ret := Info{
		CPUs:     len(cpus),
		Counters: make(map[string][]uint64),
	}

	for src.Scan() {
		items := strings.Fields(src.Text())
		var vals []uint64
		for _, item := range items[1:] {
			if v, err := strconv.ParseUint(item, 10, 64); err == nil {
				vals = append(vals, v)
			}
		}
		ret.Counters[strings.TrimSuffix(items[0], ":")] = vals
	}
	return &ret, nil
}
