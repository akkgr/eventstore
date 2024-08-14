package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/akkgr/eventstore/core"
	"github.com/akkgr/eventstore/eventstore"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)

	if err != nil {
		panic(err)
	}

	// Create a new DynamoDB client
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})

	dbc := eventstore.NewDynamoDBClient(client)

	es := eventstore.NewEventStore(dbc, dbc, core.TimerUTC{})

	events, err := es.LoadEvents("123", context.Background())
	if err != nil {
		panic(err)
	}
	prettyPrint(events)

	num := len(events)
	if num > 0 {
		err = es.Append(eventstore.Event{
			Id:      "123",
			Version: num + 1,
			Entity:  "Customer",
			Action:  "CustomerUpdated",
			Created: time.Now(),
			Data:    json.RawMessage(`{"name": "test test"}`),
		}, context.Background())
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

	events, err = es.LoadEvents("123", context.Background())

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
