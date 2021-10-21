package device

import (
	"fmt"

	"github.com/hallgren/eventsourcing"
)

// Device entity
type Device struct {
	eventsourcing.AggregateRoot
	SerialNumber string
	IP           string
	Connected    bool
}

//
// Device Events
//

// DiscoveredViaBonjour constructor event
type DiscoveredViaBonjour struct {
	IP     string
	Serial string
}

// DiscoveredViaSSDP constructor event
type DiscoveredViaSSDP struct {
	Address    string
	MacAddress string
}

// Disconnected when the device is offline
type Disconnected struct{}

// Connected when the device is online
type Connected struct{}

//
// Command errors
//

// Errors returned on failing commands
var ErrAlreadyDisconnected = fmt.Errorf("device already disconnected")
var ErrAlreadyConnected = fmt.Errorf("device already connected")

// Transitions method to build the current state of the device
func (d *Device) Transition(event eventsourcing.Event) {
	switch e := event.Data.(type) {
	case *DiscoveredViaBonjour:
		d.IP = e.IP
		d.SerialNumber = e.Serial
	case *DiscoveredViaSSDP:
		d.IP = e.Address
		d.SerialNumber = e.MacAddress
	case *Connected:
		d.Connected = true
	case *Disconnected:
		d.Connected = false
	}
}

//
// Commands
//

// Constructors
func FoundViaBonjour(ip, serial string) *Device {
	d := Device{}
	d.TrackChange(&d, &DiscoveredViaBonjour{IP: ip, Serial: serial})
	d.TrackChange(&d, &Connected{})
	return &d
}

func FoundViaSSDP(ip, serial string) *Device {
	d := Device{}
	d.TrackChange(&d, &DiscoveredViaSSDP{Address: ip, MacAddress: serial})
	d.TrackChange(&d, &Connected{})
	return &d
}

// Device commands
// NotReachable - we can't access the device
func (d *Device) NotReachable() error {
	if !d.Connected {
		return ErrAlreadyDisconnected
	}
	d.TrackChange(d, &Disconnected{})
	return nil
}

// Reachable - we can now communicate with the device
func (d *Device) Reachable() error {
	if d.Connected {
		return ErrAlreadyConnected
	}
	d.TrackChange(d, &Connected{})
	return nil
}
