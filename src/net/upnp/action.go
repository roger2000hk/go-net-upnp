// Copyright 2015 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upnp

import (
	"encoding/xml"
)

// A Action represents a icon.
type Action struct {
	XMLName      xml.Name   `xml:"action"`
	Name         string     `xml:"name"`
	ArgumentList []Argument `xml:"argumentList"`
}

// NewAction returns a new Action.
func NewAction() *Action {
	icon := &Action{}
	return icon
}
