// Copyright 2015 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
upnpdump dumps prints all devices in the local network.

        NAME
        upnpgwdump

        SYNOPSIS
        upnpdump [OPTIONS]

        DESCRIPTION
        upnpgwdump is a utility to dump SSDP messages.


        OPTIONS
        -v [0 | 1] : Enable verbose output.

        EXIT STATUS
          Return EXIT_SUCCESS or EXIT_FAILURE

        EXAMPLES
          The following is how to enable the verbose output
            upnpgwdump -v 1
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"net/upnp"
	"net/upnp/log"
)

func printGatewayDevice(n int, dev *GatewayDevice) {
	fmt.Printf("[%d] %s (%s)\n", n, dev.FriendlyName, dev.LocationURL)

	// ExternalIPAddress

	addr, err := dev.GetExternalIPAddress()
	if err == nil {
		fmt.Printf("  External IP address = %s\n", addr)
	}
}

func main() {
	// Set command line options

	verbose := flag.Int("v", 0, "Enable verbose mode [0|1]")
	flag.Usage = func() {
		cmd := strings.Split(os.Args[0], "/")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", cmd[len(cmd)-1])
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	if 0 < *verbose {
		logger := log.NewStdoutLogger(log.LoggerLevelTrace)
		log.SetSharedLogger(logger)
	}

	// Start a control point

	ctrlPoint := upnp.NewControlPoint()
	err := ctrlPoint.Start()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer ctrlPoint.Stop()

	// Search root devices

	err = ctrlPoint.SearchRootDevice()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Sleep until all search responses are received

	time.Sleep(time.Duration(ctrlPoint.SearchMX) * time.Second)

	// Print basic descriptions of found devices

	gwDevs := ctrlPoint.GetRootDevicesByType(InternetGatewayDeviceType)
	if len(gwDevs) == 0 {
		fmt.Printf("Internet gateway device is not found !!\n")
		os.Exit(0)
	}

	for n, dev := range gwDevs {
		gwDev := NewGatewayDevice(dev)
		printGatewayDevice(n, gwDev)
	}

	os.Exit(0)
}
