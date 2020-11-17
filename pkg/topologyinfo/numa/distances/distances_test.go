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

package distances

import (
	"testing"
)

type distTcase struct {
	Name string
	From int
	To   int
}

func TestDistancesBetweenNodes(t *testing.T) {
	var err error

	dists := newFakeDistances(map[int]string{
		0: "10",
	})

	val, err := dists.BetweenNodes(0, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != 10 {
		t.Errorf("unexpected distance: got %v expected %v", val, 10)
	}

	for _, tc := range []distTcase{
		{"unknown to", 0, 1},
		{"unknown from", 1, 0},
		{"negative to", 0, -1},
		{"negative from", -1, 0},
		{"both unknown", 255, 255},
		{"both negative", -1, -1},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			if _, err = dists.BetweenNodes(tc.From, tc.To); err == nil {
				t.Errorf("got distance between unknown nodes: %d, %d", tc.From, tc.To)
			}

		})
	}

}

func newFakeDistances(data map[int]string) *Distances {
	dist := Distances{
		onlineNodes: make(map[int]bool),
	}
	nodeNum := len(data)
	for nodeID, distData := range data {
		dist.onlineNodes[nodeID] = true
		nodeDist, _ := nodeDistancesFromString(nodeNum, distData)
		dist.byNode = append(dist.byNode, nodeDist)
	}
	return &dist
}
