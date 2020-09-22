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

package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fromanirh/numalign/internal/pkg/numalign"
)

func main() {
	var err error
	needSleep := false
	hours := 0 // default
	val := os.Getenv("NUMALIGN_SLEEP_HOURS")
	if val != "" {
		hours, err = strconv.Atoi(val)
		if err != nil {
			log.Fatalf("%v", err)
		}
		needSleep = true
	}

	ret := numalign.Execute()

	if needSleep {
		time.Sleep(time.Duration(hours) * time.Hour)
	}
	os.Exit(ret)
}
