package eventstore

import (
	"context"
	"errors"
	"sync"

	"github.com/akkgr/eventstore/core"
)

type AggregateId = string
type AggregateType = string
type EventName = string
type EventNumber = int
type Payload = []byte

type Event struct {
	Id      AggregateId    `json:"id" dynamodbav:"id"`
	Version EventNumber    `json:"version" dynamodbav:"version"`
	Entity  AggregateType  `json:"entity" dynamodbav:"entity"`
	Action  EventName      `json:"action" dynamodbav:"action"`
	Created core.Timestamp `json:"created" dynamodbav:"created"`
	Data    Payload        `json:"data" dynamodbav:"data"`
}

type Aggregate struct {
	Id        AggregateId    `json:"id" dynamodbav:"id"`
	Entity    AggregateType  `json:"entity" dynamodbav:"entity"`
	LastEvent Event          `json:"lastEvent" dynamodbav:"lastEvent"`
	Created   core.Timestamp `json:"created" dynamodbav:"created"`
}

type Snapshot struct {
	Id      AggregateId    `json:"id" dynamodbav:"id"`
	Version EventNumber    `json:"varsion" dynamodbav:"version"`
	Entity  AggregateType  `json:"entity" dynamodbav:"entity"`
	Created core.Timestamp `json:"created" dynamodbav:"created"`
	Data    Payload        `json:"data" dynamodbav:"data"`
}

type EventStoreAppender interface {
	Append(event Event, ctx context.Context) error
}

type EventStoreLoader interface {
	LoadEvents(id AggregateId, ctx context.Context) ([]Event, error)
}

type EventStoreReader interface {
	GetAggregate(id AggregateId, ctx context.Context) (Aggregate, error)
	GetEvents(id AggregateId, ctx context.Context) ([]Event, error)
}

type EventStoreWriter interface {
	AppendAggregate(aggregate Aggregate, ctx context.Context) error
	AppendEvent(event Event, ctx context.Context) error
	AppendSnapshot(snapshot Snapshot, ctx context.Context) error
}

type EventStore struct {
	reader EventStoreReader
	writer EventStoreWriter
	timer  core.Timer
}

func NewEventStore(reader EventStoreReader, writer EventStoreWriter, timer core.Timer) *EventStore {
	return &EventStore{
		reader: reader,
		writer: writer,
		timer:  timer,
	}
}

func (es *EventStore) Append(event Event, ctx context.Context) error {
	// get the last event
	aggregate, err := es.reader.GetAggregate(event.Id, ctx)
	if err != nil {
		return err
	}

	if aggregate.LastEvent.Version != 0 {
		// Check if the version of the last event is one less than the current event
		if aggregate.LastEvent.Version != event.Version-1 {
			return errors.New("version mismatch")
		}

		// Append the event
		err = es.writer.AppendEvent(aggregate.LastEvent, ctx)
		if err != nil {
			return err
		}
	} else {
		aggregate.Id = event.Id
		aggregate.Entity = event.Entity
		aggregate.Created = es.timer.Now()

	}

	// Update the aggregate
	aggregate.LastEvent = event
	err = es.writer.AppendAggregate(aggregate, ctx)
	return err
}

func (es *EventStore) LoadEvents(aggregateID string, ctx context.Context) ([]Event, error) {
	var wg sync.WaitGroup
	eventChan := make(chan []Event, 1)
	aggregateChan := make(chan Aggregate, 1)
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		aggregate, err := es.reader.GetAggregate(aggregateID, ctx)
		if err != nil {
			errChan <- err
			return
		}
		aggregateChan <- aggregate
	}()

	go func() {
		defer wg.Done()
		events, err := es.reader.GetEvents(aggregateID, ctx)
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

	if aggregate.LastEvent.Version == 0 {
		return events, nil
	} else {
		return append(events, aggregate.LastEvent), nil
	}
}
