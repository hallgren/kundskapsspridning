package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hallgren/eventsourcing"
	sqles "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/hallgren/kundskapsspridning/device"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// setup event store to save and get events
	es := setupEventStore()

	// create a device entity
	d := device.FoundViaBonjour("192.168.0.99", "AABBCC")
	d.AddToSite()

	spew.Dump(d)

	// repository
	//
	repository := eventsourcing.NewRepository(es, nil)
	repository.Save(d)
	fmt.Println("########")
	spew.Dump(d)

	// event subscription
	//
	// all events
	subAll := repository.SubscriberAll(func(e eventsourcing.Event) {
		fmt.Println("all", e)
	})
	subAll.Subscribe()

	// specific event
	subReason := repository.SubscriberSpecificEvent(func(e eventsourcing.Event) {
		fmt.Println("specific event", e.Reason, e)
	}, &device.RemovedFromSite{})
	subReason.Subscribe()

	// concurrency garanty
	//
	d2 := device.Device{}
	repository.Get(d.ID(), &d2)

	d2.RemoveFromSite()
	err := repository.Save(&d2)
	if err != nil {
		fmt.Println("could not store events on d2", err)
	}

	d.RemoveFromSite()
	err = repository.Save(d)
	if err != nil {
		fmt.Println("could not store events on d", err)
	}

	// global events
	//
	globalEvents, _ := es.GlobalEvents(1, 100)
	spew.Dump(globalEvents)
}

func setupEventStore() *sqles.SQL {
	serializer := eventsourcing.NewSerializer(json.Marshal, json.Unmarshal)
	serializer.RegisterTypes(&device.Device{},
		func() interface{} { return &device.DiscoveredViaBonjour{} },
		func() interface{} { return &device.AddedToSite{} },
		func() interface{} { return &device.RemovedFromSite{} },
		func() interface{} { return &device.Connected{} },
		func() interface{} { return &device.Disconnected{} },
	)

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	sqlEventStore := sqles.Open(db, *serializer)
	if err != nil {
		panic(err)
	}
	sqlEventStore.Migrate()
	return sqlEventStore
}
