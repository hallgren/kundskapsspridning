package device_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hallgren/kundskapsspridning/device"
)

func TestDiscoverViaBonjour(t *testing.T) {
	d := device.FoundViaBonjour("192.168.0.8", "AAABBB123")
	if len(d.Events()) != 2 {
		t.Fatalf("expected 2 event got %d", len(d.Events()))
	}
	if !d.Connected {
		t.Fatal("device should be connected")
	}
	spew.Dump(d.Events())
}

func TestAddToSite(t *testing.T) {
	d := device.FoundViaBonjour("192.168.0.8", "AAABBB123")
	err := d.AddToSite()
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Events()) != 3 {
		t.Fatalf("expected 3 event got %d", len(d.Events()))
	}
	if d.Events()[2].Reason != "AddedToSite" {
		t.Fatalf("last event is not AddedToSite was %s", d.Events()[1].Reason)
	}
	spew.Dump(d.Events())

	err = d.AddToSite()
	if err == nil {
		t.Fatal("should not add device to site twice")
	}
}

func TestRemoveFromSite(t *testing.T) {
	d := device.FoundViaBonjour("192.168.0.8", "AAABBB123")
	d.AddToSite()

	err := d.RemoveFromSite()
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Events()) != 4 {
		t.Fatalf("expected 4 event got %d", len(d.Events()))
	}

	if d.Events()[3].Reason != "RemovedFromSite" {
		t.Fatalf("last event should be RemovedFromSite but was: %s", d.Events()[2].Reason)
	}

	if d.PartOfSite {
		t.Fatal("the device should not be part of site")
	}
}

func TestDisconnect(t *testing.T) {
	d := device.FoundViaBonjour("192.168.0.8", "AAABBB123")
	//d.NoLongerReachable()

	if len(d.Events()) != 2 {
		t.Fatalf("expected 2 event got %d", len(d.Events()))
	}
}
