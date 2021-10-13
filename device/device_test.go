package device_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hallgren/kundskapsspridning/device"
)

func TestDiscoverViaBonjour(t *testing.T) {
	d := device.FoundViaBonjour("192.168.0.8", "AAABBB123")
	if len(d.Events()) != 1 {
		t.Fatalf("expected 1 event got %d", len(d.Events()))
	}
	spew.Dump(d.Events())
}

func TestAddToSite(t *testing.T) {
	d := device.FoundViaBonjour("192.168.0.8", "AAABBB123")
	err := d.AddToSite()
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Events()) != 2 {
		t.Fatalf("expected 2 event got %d", len(d.Events()))
	}
	if d.Events()[1].Reason != "AddedToSite" {
		t.Fatalf("2 event is not AddedToSite was %s", d.Events()[1].Reason)
	}
	spew.Dump(d.Events())

	err = d.AddToSite()
	if err == nil {
		t.Fatal("should not add device to site twice")
	}
}
