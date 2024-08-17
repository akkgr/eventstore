package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/akkgr/eventstore/core"
	"github.com/akkgr/eventstore/customer"
	"github.com/akkgr/eventstore/dynamodbstore"
	"github.com/akkgr/eventstore/eventstore"
	"github.com/akkgr/eventstore/properties"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.TODO()
	dbc := dynamodbstore.NewDynamoDBClient(ctx, true)
	createTables := flag.Bool("tbl", false, "create tables")
	flag.Parse()

	if *createTables {
		err := dynamodbstore.CreateTables(dbc)
		if err != nil {
			panic(err)
		}
	}

	myTimer := core.NewDefaultTimer()
	es := eventstore.NewEventStore(dbc, dbc, myTimer)

	// create some events
	id := uuid.New().String()
	command1 := customer.CustomerCreatedEvent{
		Name: "John Doe",
	}
	e1, _ := core.NewEvent(id, 1, "Customer", "CustomerCreated", command1)

	command2 := customer.CustomerUpdatedEvent{
		Name:   "John Doe",
		Status: "Active",
	}
	e2, _ := core.NewEvent(id, 2, "Customer", "CustomerUpdated", command2)

	command3 := customer.CustomerUpdatedEvent{
		Name:   "John Doe",
		Status: "Inactive",
		Data: properties.Data{
			Properties: map[string]properties.Property{
				"age":         properties.NewNumberProperty(decimal.NewFromInt(42)),
				"nationality": properties.NewTextProperty("American"),
				"birthday":    properties.NewDateProperty(time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC)),
			},
			Collections: map[string]properties.Collection{
				"addresses": map[string]properties.Properties{
					"home": map[string]properties.Property{
						"street":     properties.NewTextProperty("123 Main St."),
						"appartment": properties.NewTextProperty("Apt. 42"),
						"area":       properties.NewTextProperty("Springfield, IL 62701"),
					},
				},
			},
		},
	}
	e3, _ := core.NewEvent(id, 3, "Customer", "CustomerUpdated", command3)

	// publish the events
	es.Publish(e1, ctx)
	es.Publish(e2, ctx)
	es.Publish(e3, ctx)

	// read the events
	events, _ := es.LoadEvents(id, 0, ctx)
	prettyPrint(events)

	// construct the customer from the events
	c := customer.Customer{}
	for _, e := range *events {
		c.Apply(&e)
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
