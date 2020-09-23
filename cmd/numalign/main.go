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
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fromanirh/numalign/internal/pkg/numalign"
)

func main() {
	var val string
	var sleepTime time.Duration

	if _, ok := os.LookupEnv("NUMALIGN_DEBUG"); !ok {
		log.SetOutput(ioutil.Discard)
	}

	val = os.Getenv("NUMALIGN_SLEEP_HOURS")
	if val != "" {
		hours, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("%v", err)
		}
		sleepTime = time.Duration(hours) * time.Hour
	}

	log.Printf("SYS: sleep for %v after the check", sleepTime)

	R, err := numalign.NewResources()
	if err != nil {
		log.Fatalf("%v", err)
	}
	ret := numalign.Validate(R)

	if val = os.Getenv("NUMALIGN_VALIDATION_SCRIPT"); val != "" {
		code := []byte(R.MakeValidationScript())
		err := ioutil.WriteFile(val, code, 0755)
		if err != nil {
			log.Printf("SYS: validation script creation failed: %v", err)
		}
	}

	time.Sleep(sleepTime)
	os.Exit(ret)
}
