package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/akkgr/eventstore/core"
	"github.com/akkgr/eventstore/dynamodbstore"
	"github.com/akkgr/eventstore/eventstore"
)

func main() {

	createTables := flag.Bool("tbl", false, "create tables")
	flag.Parse()

	dbc := dynamodbstore.NewDynamoDBClient(context.TODO(), true)

	if *createTables {
		err := dynamodbstore.CreateTables(dbc)
		if err != nil {
			panic(err)
		}
	}

	es := eventstore.NewEventStore(dbc, dbc, core.DefaultTimer{})

	events, err := es.GetEvents("123", 0, context.Background())
	if err != nil {
		panic(err)
	}
	prettyPrint(events)

	num := len(events)
	if num > 0 {
		go func() {
			err = es.Append(eventstore.Event{
				Id:      "123",
				Version: num + 1,
				Entity:  "Customer",
				Action:  "CustomerUpdated",
				Created: time.Now(),
				Data:    json.RawMessage(`{"name": "test test"}`),
			}, context.Background())
		}()
		go func() {
			err = es.Append(eventstore.Event{
				Id:      "123",
				Version: num + 1,
				Entity:  "Customer",
				Action:  "CustomerUpdated",
				Created: time.Now(),
				Data:    json.RawMessage(`{"name": "test test"}`),
			}, context.Background())

			if err != nil {
				panic(err)
			}
		}()
	} else {
		err = es.Append(eventstore.Event{
			Id:      "123",
			Version: 1,
			Entity:  "Customer",
			Action:  "CustomerCreated",
			Created: time.Now(),
			Data:    json.RawMessage(`{"name": "test"}`),
		}, context.Background())
	}
	if err != nil {
		panic(err)
	}

	events, err = es.GetEvents("123", 0, context.Background())

	if err != nil {
		panic(err)
	}
	prettyPrint(events)
}

func prettyPrint(events []eventstore.Event) {
	b, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
