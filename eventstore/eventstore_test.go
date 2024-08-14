package eventstore_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akkgr/eventstore/eventstore"
)

// mock Timer
type mockTimer struct {
}

func (m *mockTimer) Now() time.Time {
	return time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC)
}

// mock EventStoreReader
type mockEventStoreReader struct {
}

func (m *mockEventStoreReader) GetAggregate(id eventstore.AggregateId, ctx context.Context) (eventstore.Aggregate, error) {
	switch id {
	case "1":
		return eventstore.Aggregate{}, nil
	case "2":
		return eventstore.Aggregate{}, errors.New("some error")
	default:
		return eventstore.Aggregate{
			Id:     id,
			Entity: "test",
			LastEvent: eventstore.Event{
				Id:      id,
				Entity:  "test",
				Version: 2,
				Action:  "test",
				Created: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
				Data:    []byte("test"),
			},
			Created: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
		}, nil
	}
}

func (m *mockEventStoreReader) GetEvents(id eventstore.AggregateId, ctx context.Context) ([]eventstore.Event, error) {
	switch id {
	case "1":
		return nil, errors.New("some error")
	default:
		return []eventstore.Event{
			{
				Id:      "1",
				Entity:  "test",
				Version: 1,
				Action:  "test",
				Created: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
				Data:    []byte("test"),
			},
		}, nil
	}
}

// mock EventStoreWriter
type mockEventStoreWriter struct {
}

func (m *mockEventStoreWriter) AppendAggregate(aggregate eventstore.Aggregate, ctx context.Context) error {
	switch aggregate.Id {
	case "1":
		return nil
	case "2":
		return errors.New("some error")
	default:
		return nil
	}
}

func (m *mockEventStoreWriter) AppendEvent(event eventstore.Event, ctx context.Context) error {
	switch event.Id {
	case "1":
		return nil
	case "2":
		return errors.New("some error")
	case "3":
		return errors.New("some error")
	default:
		return nil
	}
}

func (m *mockEventStoreWriter) AppendSnapshot(snapshot eventstore.Snapshot, ctx context.Context) error {
	switch snapshot.Id {
	case "1":
		return nil
	case "2":
		return errors.New("some error")
	default:
		return nil
	}
}

func CreateEventStore() *eventstore.EventStore {
	return eventstore.NewEventStore(&mockEventStoreReader{}, &mockEventStoreWriter{}, &mockTimer{})
}

func CreatedEvent(id string, version int) eventstore.Event {
	return eventstore.Event{
		Id:      id,
		Entity:  "test",
		Version: version,
		Action:  "test",
		Created: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
		Data:    []byte("test"),
	}
}

func TestAppendSuccessCreate(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("1", 2)

	// append the event
	err := es.Append(event, context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestAppendSuccessUpdate(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("4", 3)

	// append the event
	err := es.Append(event, context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestAppendFailureUpdate(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("3", 3)

	// append the event
	err := es.Append(event, context.Background())
	if err == nil {
		t.Error("expected an error")
	}
}

func TestAppendFailure(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("2", 2)

	// append the event
	err := es.Append(event, context.Background())
	if err == nil {
		t.Error("expected an error")
	}
}

func TestAppendVersionMismatch(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("3", 4)

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
	es := CreateEventStore()

	// load events
	events, err := es.LoadEvents("3", context.Background())
	if err != nil {
		t.Error(err)
	}

	if len(events) != 2 {
		t.Error("expected 1 event")
	}
}

func TestLoadEventsFailureOnGetAggregate(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// load events
	_, err := es.LoadEvents("2", context.Background())
	if err == nil {
		t.Error("expected an error")
	}
}

func TestLoadEventsFailureOnGetEvents(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// load events
	_, err := es.LoadEvents("1", context.Background())
	if err == nil {
		t.Error("expected an error")
	}
}
