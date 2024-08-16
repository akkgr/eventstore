package eventstore

import (
	"context"
	"sync"

	"github.com/akkgr/eventstore/core"
)

type EventAppender interface {
	Append(e Event, c context.Context) error
}

type EventLoader interface {
	GetEvents(id string, v int, c context.Context) ([]Event, error)
}

type EventStoreReader interface {
	GetLastEvent(id string, c context.Context) (Event, error)
	GetEvents(id string, v int, c context.Context) ([]Event, error)
}

type EventStoreWriter interface {
	AppendLastEvent(e Event, c context.Context) error
	UpdateLastEvent(e Event, c context.Context) error
	AppendEvent(e Event, c context.Context) error
}

type EventStore struct {
	reader EventStoreReader
	writer EventStoreWriter
	timer  core.Timer
}

func NewEventStore(r EventStoreReader, w EventStoreWriter, t core.Timer) *EventStore {
	return &EventStore{
		reader: r,
		writer: w,
		timer:  t,
	}
}

func (es *EventStore) Append(e Event, c context.Context) error {
	a, err := es.reader.GetLastEvent(e.Id, c)
	if err != nil {
		return err
	}

	if a.Version != e.Version-1 {
		return InvalidVersion{}
	}

	if a.Version > 0 {
		err = es.writer.AppendEvent(a, c)
		if err != nil {
			return err
		}
		err = es.writer.UpdateLastEvent(e, c)
		return err
	} else {
		err = es.writer.AppendLastEvent(e, c)
		return err
	}
}

// LoadEvents loads all events for an aggregate starting from a specific version.
// If the version is 0, it will return all events.
// If the version is greater than the current version, it will return an error.
// id is the aggregate id.
// v is the version to start from.
// c is the context.
func (es *EventStore) GetEvents(id string, v int, c context.Context) ([]Event, error) {
	var wg sync.WaitGroup
	eventChan := make(chan []Event, 1)
	aggregateChan := make(chan Event, 1)
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		aggregate, err := es.reader.GetLastEvent(id, c)
		if err != nil {
			errChan <- err
			return
		}
		aggregateChan <- aggregate
	}()

	go func() {
		defer wg.Done()
		events, err := es.reader.GetEvents(id, v, c)
		if err != nil {
			errChan <- err
			return
		}
		eventChan <- events
	}()

	wg.Wait()
	close(eventChan)
	close(aggregateChan)
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	aggregate := <-aggregateChan
	events := <-eventChan

	if aggregate.Version < v {
		return nil, InvalidVersion{}
	}

	// return the events up the persisted aggregate version
	// due to concurrency errors, there might be more events than the last version
	// these events are ignored and will be overridden in future appends
	if aggregate.Version > 0 {
		events = events[:aggregate.Version-1]
	}

	return append(events, aggregate), nil
}
