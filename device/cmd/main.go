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

	// repository is the kit between the entity and the event store
	repository := eventsourcing.NewRepository(es, nil)

	// create a device entity
	d := device.FoundViaBonjour("192.168.0.99", "AABBCC")
	spew.Dump(d.Events())

	// store d events in the repository
	repository.Save(d)
	fmt.Println("######")
	spew.Dump(d.Events())

	/*
		// fetch the same device events from the repository and build the entity
		d2 := device.Device{}
		repository.Get(d.ID(), &d2)

		fmt.Println("########")
		spew.Dump(d2)

			// optimistic concurrency
			//
			d2.NotReachable()
			err := repository.Save(&d2)
			if err != nil {
				fmt.Println("could not store events on d2", err)
			}

			d.NotReachable()
			err = repository.Save(d)
			if err != nil {
				fmt.Println("could not store events on d", err)
			}

				// event subscription
				//
				// all events
				subAll := repository.SubscriberAll(func(e eventsourcing.Event) {
					fmt.Println("all", e.Reason, e)
				})

				// specific event
				subReason := repository.SubscriberSpecificEvent(func(e eventsourcing.Event) {
					fmt.Println("specific event", e.Reason, e)
				}, &device.Connected{}, &device.Disconnected)

				subAll.Subscribe()
				subReason.Subscribe()
					// global events
					//
					globalEvents, _ := es.GlobalEvents(1, 100)
					spew.Dump(globalEvents)
	*/
}

func setupEventStore() *sqles.SQL {
	serializer := eventsourcing.NewSerializer(json.Marshal, json.Unmarshal)
	serializer.RegisterTypes(&device.Device{},
		func() interface{} { return &device.DiscoveredViaBonjour{} },
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
