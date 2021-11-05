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

func TestDisconnectConnect(t *testing.T) {
	d := device.FoundViaBonjour("192.168.0.8", "AAABBB123")

	// Device disconnected
	err := d.NotReachable()
	if err != nil {
		t.Fatal(err)
	}

	if len(d.Events()) != 3 {
		t.Fatalf("expected 3 event got %d", len(d.Events()))
	}

	if d.Events()[2].Reason() != "Disconnected" {
		t.Fatalf("expected last event to be Disconnected but was %s", d.Events()[2].Reason())
	}

	// device is connected
	err = d.Reachable()
	if err != nil {
		t.Fatal(err)
	}

	// last event should be a connected event
	if d.Events()[3].Reason() != "Connected" {
		t.Fatalf("expected last event to be Connected but was %s", d.Events()[3].Reason())
	}

	// test that current state is updated
	if !d.Connected {
		t.Fatal("device is disconnected")
	}
}
