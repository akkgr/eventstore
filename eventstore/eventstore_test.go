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

func (m *mockEventStoreReader) GetLastEvent(id string, c context.Context) (eventstore.Event, error) {
	switch id {
	case "no events":
		return eventstore.Event{}, nil
	case "call error":
		return eventstore.Event{}, errors.New("some error")
	default:
		return eventstore.Event{
			Id:      id,
			Version: 2,
			Entity:  "test",
			Action:  "test",
			Created: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
			Data:    []byte("test"),
		}, nil
	}
}

func (m *mockEventStoreReader) GetEvents(id string, v int, c context.Context) ([]eventstore.Event, error) {
	switch id {
	case "no events":
		return nil, nil
	case "call error":
		return nil, errors.New("some error")
	default:
		return []eventstore.Event{
			{
				Id:      id,
				Version: 1,
				Entity:  "test",
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

func (m *mockEventStoreWriter) AppendLastEvent(e eventstore.Event, c context.Context) error {
	switch e.Id {
	case "call error":
		return errors.New("some error")
	default:
		return nil
	}
}

func (m *mockEventStoreWriter) UpdateLastEvent(e eventstore.Event, c context.Context) error {
	switch e.Id {
	case "call error":
		return errors.New("some error")
	default:
		return nil
	}
}

func (m *mockEventStoreWriter) AppendEvent(e eventstore.Event, c context.Context) error {
	switch e.Id {
	case "call error":
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
		Version: version,
		Entity:  "test",
		Action:  "test",
		Created: time.Date(2000, 3, 15, 0, 0, 0, 0, time.UTC),
		Data:    []byte("test"),
	}
}

func TestFirstEventSuccess(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("no events", 1)

	// append the event
	err := es.Append(event, context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestSecondEventSuccess(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("id", 3)

	// append the event
	err := es.Append(event, context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestAppendFailure(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// create a new event
	event := CreatedEvent("call error", 2)

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
	event := CreatedEvent("id", 4)

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
	events, err := es.GetEvents("id", 0, context.Background())
	if err != nil {
		t.Error(err)
	}

	if len(events) != 2 {
		t.Error("expected 1 event")
	}
}

func TestLoadEventsFailure(t *testing.T) {
	// create a new event store
	es := CreateEventStore()

	// load events
	_, err := es.GetEvents("call error", 0, context.Background())
	if err == nil {
		t.Error("expected an error")
	}
}
