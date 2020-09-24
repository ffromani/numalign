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

	flag "github.com/spf13/pflag"

	"github.com/fromanirh/numalign/internal/pkg/numalign"
)

func main() {
	var sleepTime time.Duration

	var sleepHoursParam = flag.StringP("sleep-hours", "S", "", "sleep hours once done.")
	var scriptPathParam = flag.StringP("script-path", "P", "", "save test script to this path.")
	flag.Parse()

	if _, ok := os.LookupEnv("NUMALIGN_DEBUG"); !ok {
		log.SetOutput(ioutil.Discard)
	}

	sleepHours := *sleepHoursParam
	if sleepHours == "" {
		os.Getenv("NUMALIGN_SLEEP_HOURS")
	}
	if sleepHours != "" {
		hours, err := strconv.Atoi(sleepHours)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if hours > 0 {
			sleepTime = time.Duration(hours) * time.Hour
		}
	}

	log.Printf("SYS: sleep for %v after the check", sleepTime)

	R, err := numalign.NewResources(flag.Args())
	if err != nil {
		log.Fatalf("%v", err)
	}
	ret := numalign.Validate(R)

	scriptPath := *scriptPathParam
	if scriptPath == "" {
		scriptPath = os.Getenv("NUMALIGN_VALIDATION_SCRIPT")
	}

	if scriptPath != "" {
		code := []byte(R.MakeValidationScript())
		err := ioutil.WriteFile(scriptPath, code, 0755)
		if err != nil {
			log.Printf("SYS: validation script creation failed: %v", err)
		}
	}

	time.Sleep(sleepTime)
	os.Exit(ret)
}
