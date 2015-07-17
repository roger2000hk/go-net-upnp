// Copyright 2015 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
upnpdump dumps prints all devices in the local network.

        NAME
        upnpdump

        SYNOPSIS
        upnpdump [OPTIONS]

        DESCRIPTION
        upnpdump is a utility to dump SSDP messages.


        OPTIONS
        -v [0 | 1] : Enable verbose output.

        EXIT STATUS
          Return EXIT_SUCCESS or EXIT_FAILURE

        EXAMPLES
          The following is how to enable the verbose output
            upnpdump -v 1
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"net/upnp"
	"net/upnp/log"
)

func main() {
	// Set command line options

	verbose := flag.Int("v", 0, "Output log level.")
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

	if len(ctrlPoint.GetRootDevices()) == 0 {
		fmt.Printf("Device not found !!")
		os.Exit(0)
	}

	for n, dev := range ctrlPoint.GetRootDevices() {
		devURL := dev.LocationURL

		presentationURL := dev.PresentationURL
		if 0 < len(presentationURL) {
			url, err := dev.GetAbsoluteURL(presentationURL)
			if err == nil {
				devURL = url.String()
			}
		}

		fmt.Printf("[%d] '%s', '%s', %s\n", n, dev.FriendlyName, dev.DeviceType, devURL)
	}

	os.Exit(0)
}
