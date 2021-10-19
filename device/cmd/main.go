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
	err = sqlEventStore.Migrate()
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}

	// Create a repository to handle event sourced
	d := device.FoundViaBonjour("192.168.0.99", "AABBCC")
	d.AddToSite()

	spew.Dump(d)

	//repository := eventsourcing.NewRepository(sqlEventStore, nil)
	//repository.Save(d)

	//fmt.Println("########")
	//d2 := device.Device{}
	//repository.Get(d.ID(), &d2)
	//spew.Dump(d2)

}
