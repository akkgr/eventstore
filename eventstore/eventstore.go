package eventstore

import (
	"context"
	"sync"

	"github.com/akkgr/eventstore/core"
)

type EventPublisher interface {
	Publish(e *core.Event, c context.Context) error
}

type EventLoader interface {
	LoadEvents(id string, v int, c context.Context) (*[]core.Event, error)
}

type EventStoreReader interface {
	GetLastEvent(id string, c context.Context) (*core.Event, error)
	GetEvents(id string, v int, c context.Context) (*[]core.Event, error)
}

type EventStoreWriter interface {
	AppendLastEvent(e *core.Event, c context.Context) error
	UpdateLastEvent(e *core.Event, c context.Context) error
	AppendEvent(e *core.Event, c context.Context) error
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

func (es *EventStore) Publish(e *core.Event, c context.Context) error {
	a, err := es.reader.GetLastEvent(e.Id, c)
	if err != nil {
		return err
	}

	if a.Version != e.Version-1 {
		return core.InvalidVersion{}
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
func (es *EventStore) LoadEvents(id string, v int, c context.Context) (*[]core.Event, error) {
	var wg sync.WaitGroup
	eventChan := make(chan *[]core.Event, 1)
	aggregateChan := make(chan *core.Event, 1)
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

	if aggregate.Version == 0 {
		return nil, core.EventsNotFound{}
	}

	if aggregate.Version < v {
		return nil, core.InvalidVersion{}
	}

	all := *events
	last := *aggregate

	// due to concurrency errors, there might be more events than the last version
	// these events are ignored and will be overridden in future appends
	if len(all) > 0 && all[len(all)-1].Version > aggregate.Version {
		all = all[:aggregate.Version]
		events = &all
	}

	result := append(all, last)

	return &result, nil
}
