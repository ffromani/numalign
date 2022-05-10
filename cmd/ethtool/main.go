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
 * Copyright 2022 Red Hat, Inc.
 */

package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/fromanirh/ethtool"
)

func main() {
	var showChannels = flag.BoolP("show-channels", "l", false, "show channel info for device(s)")
	flag.Parse()

	ifaces := flag.Args()
	cli, err := ethtool.New()
	mustNotFail(err, 1)
	defer cli.Close()

	if len(ifaces) == 0 {
		if *showChannels {
			cis, err := cli.ChannelInfos()
			mustNotFail(err, 2)
			for _, ci := range cis {
				showChannelInfo(ci)
			}
		}
		os.Exit(0)
	}

	for _, iface := range ifaces {
		if *showChannels {
			ci, err := cli.ChannelInfo(ethtool.Interface{Name: iface})
			mustNotFail(err, 2)
			showChannelInfo(ci)
		}

	}
}

func mustNotFail(err error, code int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(code)
	}
}

func showChannelInfo(ci *ethtool.ChannelInfo) {
	fmt.Printf("Channel parameters for %s:\n", ci.Interface.Name)
	fmt.Printf("Pre-set maximums:\n")
	fmt.Printf("RX:		%d\n", ci.MaxRx)
	fmt.Printf("TX:		%d\n", ci.MaxTx)
	fmt.Printf("Other:		%d\n", ci.MaxOther)
	fmt.Printf("Combined:	%d\n", ci.MaxCombined)
	fmt.Printf("Current hardware settings:\n")
	fmt.Printf("RX:		%d\n", ci.Rx)
	fmt.Printf("TX:		%d\n", ci.Tx)
	fmt.Printf("Other:		%d\n", ci.Other)
	fmt.Printf("Combined:	%d\n", ci.Combined)
}
