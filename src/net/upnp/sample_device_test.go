// Copyright 2015 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upnp

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"testing"

	"net/upnp/control"
)

const (
	errorTestDeviceInvalidURL           = "invalid url %s = '%s', expected : '%s'"
	errorTestDeviceInvalidStatusCode    = "invalid status code (%s) = [%d] : expected : [%d]"
	errorTestDeviceInvalidPortRange     = "invalid port range = [%d] : expected : [%d]~[%d]"
	errorTestDeviceInvalidParentObject  = "invalid parent object %p = '%p', expected : '%p'"
	errorTestDeviceInvalidArgumentValue = "invalid argument value %s = '%s', expected : '%s'"
	errorTestDeviceInvalidArgumentDir   = "invalid argument direction %s = %d, expected : %d"
)

const (
	SetTarget      = "SetTarget"
	GetTarget      = "GetTarget"
	NewTargetValue = "newTargetValue"
	RetTargetValue = "RetTargetValue"
)

type sampleDevice struct {
	*Device
	Target string
}

func NewSampleDevice() (*sampleDevice, error) {
	dev, err := NewDeviceFromDescription(binaryLightDeviceDescription)
	if err != nil {
		return nil, err
	}

	service, err := dev.GetServiceByType("urn:schemas-upnp-org:service:SwitchPower:1")
	if err != nil {
		return nil, err
	}

	err = service.LoadDescriptionBytes([]byte(switchPowerServiceDescription))
	if err != nil {
		return nil, err
	}

	sampleDev := &sampleDevice{Device: dev}
	sampleDev.ActionListener = sampleDev

	return sampleDev, nil
}

func (self *sampleDevice) GetSwitchPowerService() (*Service, error) {
	return self.GetServiceByType("urn:schemas-upnp-org:service:SwitchPower:1")
}

func (self *sampleDevice) GetSwitchPowerSetTargetAction() (*Action, error) {
	service, err := self.GetSwitchPowerService()
	if err != nil {
		return nil, err
	}
	return service.GetActionByName(SetTarget)
}

func (self *sampleDevice) GetSwitchPowerGetTargetAction() (*Action, error) {
	service, err := self.GetSwitchPowerService()
	if err != nil {
		return nil, err
	}
	return service.GetActionByName(GetTarget)
}

func (self *sampleDevice) ActionRequestReceived(action *Action) *control.UPnPError {
	switch action.Name {
	case SetTarget:
		arg, err := action.GetArgumentByName(NewTargetValue)
		if err == nil {
			self.Target = arg.Value
		}
		return nil
	case GetTarget:
		arg, err := action.GetArgumentByName(RetTargetValue)
		if err == nil {
			arg.Value = self.Target
		}
		return nil
	}

	return control.NewUPnPErrorFromCode(control.ErrorOptionalActionNotImplemented)
}

func TestSampleDeviceDescription(t *testing.T) {
	dev, err := NewSampleDevice()

	if err != nil {
		t.Error(err)
	}

	// check service

	service, err := dev.GetServiceByType("urn:schemas-upnp-org:service:SwitchPower:1")
	if err != nil {
		t.Error(err)
	}

	if service.ParentDevice != dev.Device {
		t.Errorf(errorTestDeviceInvalidParentObject, service, service.ParentDevice, dev.Device)
	}

	service, err = dev.GetServiceById("urn:upnp-org:serviceId:SwitchPower.1")
	if err != nil {
		t.Error(err)
	}

	if service.ParentDevice != dev.Device {
		t.Errorf(errorTestDeviceInvalidParentObject, service, service.ParentDevice, dev.Device)
	}

	// check actions

	actionNames := []string{"SetTarget", "GetTarget", "GetStatus"}
	for _, name := range actionNames {
		action, err := service.GetActionByName(name)
		if err != nil {
			t.Error(err)
		}
		if action.ParentService != service {
			t.Errorf(errorTestDeviceInvalidParentObject, action, action.ParentService, service)
		}
	}

	// check argumengs (SetTarget)

	action, err := service.GetActionByName("SetTarget")
	if err == nil {
		argNames := []string{"newTargetValue"}
		argDirs := []int{InDirection}
		for n, name := range argNames {
			arg, err := action.GetArgumentByName(name)
			if err != nil {
				t.Error(err)
			}

			argDir := arg.GetDirection()
			if argDir != argDirs[n] {
				t.Errorf(errorTestDeviceInvalidArgumentDir, name, argDir, argDirs[n])
			}

			// check parent service

			if arg.ParentAction != action {
				t.Errorf(errorTestDeviceInvalidParentObject, arg, arg.ParentAction, action)
			}

			// check setter and getter

			value := fmt.Sprintf("%d", rand.Int())
			err = arg.SetString(value)
			if err != nil {
				t.Error(err)
			}
			argValue, err := arg.GetString()
			if err != nil {
				t.Error(err)
			}
			if value != argValue {
				t.Errorf(errorTestDeviceInvalidArgumentValue, name, argValue, value)
			}
		}
	} else {
		t.Error(err)
	}

	// start device

	err = dev.Start()
	if err != nil {
		t.Error(err)
	}

	// check service

	checkServiceURLs := func(dev *sampleDevice, serviceType string, urls []string) {
		service, err := dev.GetServiceByType(serviceType)
		if err != nil {
			t.Error(err)
		}

		expectURL := urls[0]
		if len(service.SCPDURL) <= 0 || service.SCPDURL != expectURL {
			t.Errorf(errorTestDeviceInvalidURL, "SCPDURL", service.SCPDURL, expectURL)
		}

		expectURL = urls[1]
		if len(service.ControlURL) <= 0 || service.ControlURL != expectURL {
			t.Errorf(errorTestDeviceInvalidURL, "ControlURL", service.ControlURL, expectURL)
		}

		expectURL = urls[2]
		if len(service.EventSubURL) <= 0 || service.EventSubURL != expectURL {
			t.Errorf(errorTestDeviceInvalidURL, "EventSubURL", service.EventSubURL, expectURL)
		}
	}

	urls := []string{
		"/service/scpd/SwitchPower.xml",
		"/service/control/SwitchPower",
		"/service/event/SwitchPower"}
	checkServiceURLs(dev, "urn:schemas-upnp-org:service:SwitchPower:1", urls)

	// stop device
	err = dev.Stop()
	if err != nil {
		t.Error(err)
	}
}

const binaryLightDeviceDescription = xml.Header +
	"<root>" +
	"  <device>" +
	"    <serviceList>" +
	"      <service>" +
	"        <serviceType>urn:schemas-upnp-org:service:SwitchPower:1</serviceType>" +
	"        <serviceId>urn:upnp-org:serviceId:SwitchPower.1</serviceId>" +
	"      </service>" +
	"    </serviceList>" +
	"  </device>" +
	"</root>"

const switchPowerServiceDescription = xml.Header +
	"<scpd>" +
	"  <serviceStateTable>" +
	"    <stateVariable>" +
	"      <name>Target</name>" +
	"      <sendEventsAttribute>no</sendEventsAttribute> " +
	"      <dataType>boolean</dataType>" +
	"      <defaultValue>0</defaultValue>" +
	"    </stateVariable>" +
	"    <stateVariable>" +
	"      <name>Status</name>" +
	"      <dataType>boolean</dataType>" +
	"      <defaultValue>0</defaultValue>" +
	"    </stateVariable>" +
	"  </serviceStateTable>" +
	"  <actionList>" +
	"    <action>" +
	"    <name>SetTarget</name>" +
	"      <argumentList>" +
	"        <argument>" +
	"          <name>newTargetValue</name>" +
	"          <direction>in</direction>" +
	"          <relatedStateVariable>Target</relatedStateVariable>" +
	"        </argument>" +
	"      </argumentList>" +
	"    </action>" +
	"    <action>" +
	"    <name>GetTarget</name>" +
	"      <argumentList>" +
	"        <argument>" +
	"          <name>RetTargetValue</name>" +
	"          <direction>out</direction>" +
	"          <relatedStateVariable>Target</relatedStateVariable>" +
	"        </argument>" +
	"      </argumentList>" +
	"    </action>" +
	"    <action>" +
	"    <name>GetStatus</name>" +
	"      <argumentList>" +
	"        <argument>" +
	"          <name>ResultStatus</name>" +
	"          <direction>out</direction>" +
	"          <relatedStateVariable>Status</relatedStateVariable>" +
	"        </argument>" +
	"      </argumentList>" +
	"    </action>" +
	"  </actionList>" +
	"</scpd>"
