package device

import (
	"fmt"

	"github.com/hallgren/eventsourcing"
)

// Device entity
type Device struct {
	eventsourcing.AggregateRoot
	Serial     string
	Connected  bool
	PartOfSite bool
	IP         string
}

//
// Device Events
//

// DiscoveredViaBonjour constructor event
type DiscoveredViaBonjour struct {
	IP     string
	Serial string
}

// AddedToSite when the device is part of the site
type AddedToSite struct{}

// RemovedFromSite when the device is no longer part of the site
type RemovedFromSite struct{}

// Disconnected when the device is offline
type Disconnected struct{}

// Connected when the device is online
type Connected struct{}

//
// Device errors
//

// Errors returned on failing commands
var ErrAlreadyPartOfSite = fmt.Errorf("device is already part of site")
var ErrNotPartOfSite = fmt.Errorf("device is not part of site")
var ErrAlreadyDisconnected = fmt.Errorf("device already disconnected")
var ErrAlreadyConnected = fmt.Errorf("device already connected")

// Transitions method to build the current state of the device
func (d *Device) Transition(event eventsourcing.Event) {
	switch e := event.Data.(type) {
	case *DiscoveredViaBonjour:
		d.IP = e.IP
		d.Serial = e.Serial
	case *Connected:
		d.Connected = true
	case *Disconnected:
		d.Connected = false
	case *AddedToSite:
		d.PartOfSite = true
	case *RemovedFromSite:
		d.PartOfSite = false
	}
}

//
// Commands
//

// Constructor
func FoundViaBonjour(ip, serial string) *Device {
	d := Device{}
	d.TrackChange(&d, &DiscoveredViaBonjour{IP: ip, Serial: serial})
	d.TrackChange(&d, &Connected{})
	return &d
}

func (d *Device) AddToSite() error {
	if d.PartOfSite {
		return ErrAlreadyPartOfSite
	}
	d.TrackChange(d, &AddedToSite{})
	return nil
}

func (d *Device) RemoveFromSite() error {
	if !d.PartOfSite {
		return ErrNotPartOfSite
	}
	d.TrackChange(d, &RemovedFromSite{})
	return nil
}

func (d *Device) Disconnect() error {
	if !d.Connected {
		return ErrAlreadyDisconnected
	}
	d.TrackChange(d, &Disconnected{})
	return nil
}

func (d *Device) Connect() error {
	if d.Connected {
		return ErrAlreadyConnected
	}
	d.TrackChange(d, &Connected{})
	return nil
}
