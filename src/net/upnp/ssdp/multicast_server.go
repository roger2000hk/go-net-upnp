// Copyright 2015 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssdp

import (
	"net/upnp/log"
)

// A SSDPListener represents a listener for MulticastServer.
type MulticastListener interface {
	DeviceNotifyReceived(ssdpReq *Request)
	DeviceSearchReceived(ssdpReq *Request)
}

// A MulticastServer represents a packet of SSDP.
type MulticastServer struct {
	Socket   *HTTPMUSocket
	Listener MulticastListener
}

// NewMulticastServer returns a new MulticastServer.
func NewMulticastServer() *MulticastServer {
	ssdpPkt := &MulticastServer{}
	ssdpPkt.Socket = NewHTTPMUSocket()
	ssdpPkt.Listener = nil
	return ssdpPkt
}

// Start starts this server.
func (self *MulticastServer) Start() error {
	err := self.Socket.Bind()
	if err != nil {
		return err
	}
	go handleMulticastConnection(self)
	return nil
}

// Stop stops this server.
func (self *MulticastServer) Stop() error {
	err := self.Socket.Close()
	if err != nil {
		return err
	}
	return nil
}

func handleMulticastConnection(self *MulticastServer) {
	for {
		ssdpPkt, err := self.Socket.Read()
		if err != nil {
			log.Error(err)
			break
		}

		if len(ssdpPkt.Bytes) <= 0 {
			continue
		}

		if self.Listener != nil {
			ssdpReq, _ := NewRequestFromPacket(ssdpPkt)
			switch {
			case ssdpReq.IsNotifyRequest():
				self.Listener.DeviceNotifyReceived(ssdpReq)
			case ssdpReq.IsSearchRequest():
				self.Listener.DeviceSearchReceived(ssdpReq)
			}
		}
	}
}