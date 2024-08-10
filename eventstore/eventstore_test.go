package eventstore_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akkgr/eventstore/eventstore"
)

// mock EventStoreReader
type mockEventStoreReader struct {
}

func (m *mockEventStoreReader) GetAggregate(aggregateID string, ctx context.Context) (eventstore.Aggregate, error) {
	switch aggregateID {
	case "1":
		return eventstore.Aggregate{}, nil
	case "2":
		return eventstore.Aggregate{}, errors.New("some error")
	default:
		return eventstore.Aggregate{
			AggregateID:   aggregateID,
			AggregateType: "test",
			LastEvent: eventstore.Event{
				AggregateID:   aggregateID,
				AggregateType: "test",
				EventNumber:   2,
				EventName:     "test",
				Created:       time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
				Data:          []byte("test"),
			},
			Created: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
			Updated: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
		}, nil
	}
}

func (m *mockEventStoreReader) GetEvents(aggregateID string, ctx context.Context) ([]eventstore.Event, error) {
	switch aggregateID {
	case "1":
		return nil, errors.New("some error")
	default:
		return []eventstore.Event{
			{
				AggregateID:   "1",
				AggregateType: "test",
				EventNumber:   1,
				EventName:     "test",
				Created:       time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
				Data:          []byte("test"),
			},
		}, nil
	}
}

// mock EventStoreWriter
type mockEventStoreWriter struct {
}

func (m *mockEventStoreWriter) AppendAggregate(aggregate eventstore.Aggregate, ctx context.Context) error {
	if aggregate.AggregateID == "1" {
		return nil
	} else {
		return errors.New("some error")
	}
}

func (m *mockEventStoreWriter) AppendEvent(event eventstore.Event, ctx context.Context) error {
	if event.AggregateID == "1" {
		return nil
	} else {
		return errors.New("some error")
	}
}

func TestAppendSuccess(t *testing.T) {
	// create a new event store
	es := eventstore.NewEventStore(&mockEventStoreReader{}, &mockEventStoreWriter{})

	// create a new event
	event := eventstore.Event{
		AggregateID:   "1",
		AggregateType: "test",
		EventNumber:   2,
		EventName:     "test",
		Created:       time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
		Data:          []byte("test"),
	}

	// append the event
	err := es.Append(event, context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestAppendFailure(t *testing.T) {
	// create a new event store
	es := eventstore.NewEventStore(&mockEventStoreReader{}, &mockEventStoreWriter{})

	// create a new event
	event := eventstore.Event{
		AggregateID:   "2",
		AggregateType: "test",
		EventNumber:   2,
		EventName:     "test",
		Created:       time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
		Data:          []byte("test"),
	}

	// append the event
	err := es.Append(event, context.Background())
	if err == nil {
		t.Error("expected an error")
	}
}

func TestAppendVersionMismatch(t *testing.T) {
	// create a new event store
	es := eventstore.NewEventStore(&mockEventStoreReader{}, &mockEventStoreWriter{})

	// create a new event
	event := eventstore.Event{
		AggregateID:   "3",
		AggregateType: "test",
		EventNumber:   4,
		EventName:     "test",
		Created:       time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
		Data:          []byte("test"),
	}

	// append the event
	err := es.Append(event, context.Background())
	if err == nil {
		t.Error("expected an error")
	}
	if condition := err.Error(); condition != "version mismatch" {
		t.Errorf("expected version mismatch, got %s", condition)
	}
}

func TestLoadEventsSuccess(t *testing.T) {
	// create a new event store
	es := eventstore.NewEventStore(&mockEventStoreReader{}, &mockEventStoreWriter{})

	// load events
	events, err := es.LoadEvents("3", context.Background())
	if err != nil {
		t.Error(err)
	}

	if len(events) != 2 {
		t.Error("expected 1 event")
	}
}

func TestLoadEventsFailure(t *testing.T) {
	// create a new event store
	es := eventstore.NewEventStore(&mockEventStoreReader{}, &mockEventStoreWriter{})

	// load events
	_, err := es.LoadEvents("1", context.Background())
	if err == nil {
		t.Error("expected an error")
	}
}
