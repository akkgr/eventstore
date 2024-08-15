package eventstore_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/akkgr/eventstore/core"
	"github.com/akkgr/eventstore/dynamodbstore"
	"github.com/akkgr/eventstore/eventstore"
)

func BenchmarkStoreLoadEvents(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	dbc := dynamodbstore.NewDynamoDBClient(context.TODO(), true)
	es := eventstore.NewEventStore(dbc, dbc, core.DefaultTimer{})

	for i := 0; i < b.N; i++ {
		_, err := es.LoadEvents("123", 0, context.Background())
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkStoreAppendEvents(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	dbc := dynamodbstore.NewDynamoDBClient(context.TODO(), true)
	es := eventstore.NewEventStore(dbc, dbc, core.DefaultTimer{})

	for i := 0; i < b.N; i++ {
		_ = es.Append(eventstore.Event{
			Id:      "123",
			Version: i + 1,
			Entity:  "Customer",
			Action:  "CustomerUpdated",
			Created: time.Now(),
			Data:    json.RawMessage(`{"name": "test test"}`),
		}, context.Background())
	}
}
