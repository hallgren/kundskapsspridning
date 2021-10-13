package device

import (
	"fmt"

	"github.com/hallgren/eventsourcing"
)

type Device struct {
	eventsourcing.AggregateRoot
	Serial     string
	Connected  bool
	PartOfSite bool
	IP         string
}

type DiscoveredViaBonjour struct {
	IP     string
	Serial string
}

type AddedToSite struct{}

var ErrAlreadyPartOfSite = fmt.Errorf("device is already part of site")

func (d *Device) Transition(event eventsourcing.Event) {
	switch e := event.Data.(type) {
	case *DiscoveredViaBonjour:
		d.IP = e.IP
		d.Serial = e.Serial
	case *AddedToSite:
		d.PartOfSite = true
	}
}

func FoundViaBonjour(ip, serial string) *Device {
	d := Device{}
	d.TrackChange(&d, &DiscoveredViaBonjour{IP: ip, Serial: serial})
	return &d
}

func (d *Device) AddToSite() error {
	if d.PartOfSite {
		return ErrAlreadyPartOfSite
	}
	d.TrackChange(d, &AddedToSite{})
	return nil
}
