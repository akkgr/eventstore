package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/akkgr/eventstore/core"
	"github.com/akkgr/eventstore/customer"
	"github.com/akkgr/eventstore/dynamodbstore"
	"github.com/akkgr/eventstore/eventstore"
	"github.com/google/uuid"
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

	es := eventstore.NewEventStore(dbc, dbc, core.NewDefaultTimer())

	// create some events
	id := uuid.New().String()
	e1, _ := core.NewEvent(id, 1, "Customer", "CustomerCreated", map[string]interface{}{"Name": "John Doe"})
	e2, _ := core.NewEvent(id, 2, "Customer", "CustomerUpdated", map[string]interface{}{"Name": "John Doe", "Status": "Active"})
	e3, _ := core.NewEvent(id, 3, "Customer", "CustomerUpdated", map[string]interface{}{"Name": "John Doe", "Status": "Inactive"})

	// publish the events
	es.Publish(e1, context.Background())
	es.Publish(e2, context.Background())
	es.Publish(e3, context.Background())

	// read the events
	events, _ := es.LoadEvents(id, 0, context.Background())
	prettyPrint(events)

	// construct the customer from the events
	c := customer.Customer{}
	for _, e := range *events {
		c.Apply(e)
	}
	prettyPrint(c)
}

func prettyPrint(obj any) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
