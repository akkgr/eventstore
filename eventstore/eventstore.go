package eventstore

import (
	"context"
	"errors"
	"time"
)

type Aggregate struct {
	AggregateID   string    `json:"aggregateId" dynamodbav:"aggregateId"`
	AggregateType string    `json:"aggregateType" dynamodbav:"aggregateType"`
	LastEvent     Event     `json:"lastEvent" dynamodbav:"lastEvent"`
	Created       time.Time `json:"created" dynamodbav:"created"`
	Updated       time.Time `json:"updated" dynamodbav:"updated"`
}

type Event struct {
	AggregateID   string    `json:"aggregateId" dynamodbav:"aggregateId"`
	EventNumber   int       `json:"eventNumber" dynamodbav:"eventNumber"`
	AggregateType string    `json:"aggregateType" dynamodbav:"aggregateType"`
	EventName     string    `json:"eventName" dynamodbav:"eventName"`
	Created       time.Time `json:"created" dynamodbav:"created"`
	Data          []byte    `json:"data" dynamodbav:"data"`
}

type EventStoreAppender interface {
	Append(event Event, ctx context.Context) error
}

type EventStoreLoader interface {
	LoadEvents(aggregateID string, ctx context.Context) ([]Event, error)
}

type EventStoreReader interface {
	GetAggregate(aggregateID string, ctx context.Context) (Aggregate, error)
	GetEvents(aggregateID string, ctx context.Context) ([]Event, error)
}

type EventStoreWriter interface {
	AppendAggregate(aggregate Aggregate, ctx context.Context) error
	AppendEvent(event Event, ctx context.Context) error
}

type EventStore struct {
	reader EventStoreReader
	writer EventStoreWriter
}

func NewEventStore(reader EventStoreReader, writer EventStoreWriter) *EventStore {
	return &EventStore{
		reader: reader,
		writer: writer,
	}
}

func (es *EventStore) Append(event Event, ctx context.Context) error {
	// get the last event
	aggregate, err := es.reader.GetAggregate(event.AggregateID, ctx)
	if err != nil {
		return err
	}

	if aggregate.LastEvent.EventNumber != 0 {
		// Check if the version of the last event is one less than the current event
		if aggregate.LastEvent.EventNumber != event.EventNumber-1 {
			return errors.New("version mismatch")
		}

		// Append the event
		err = es.writer.AppendEvent(aggregate.LastEvent, ctx)
		if err != nil {
			return err
		}
		aggregate.Updated = time.Now()

	} else {
		aggregate.AggregateID = event.AggregateID
		aggregate.AggregateType = event.AggregateType
		aggregate.Created = time.Now()

	}

	// Update the aggregate
	aggregate.LastEvent = event
	err = es.writer.AppendAggregate(aggregate, ctx)
	return err
}

func (es *EventStore) LoadEvents(aggregateID string, ctx context.Context) ([]Event, error) {
	aggregate, err := es.reader.GetAggregate(aggregateID, ctx)
	if err != nil {
		return nil, err
	}

	events, err := es.reader.GetEvents(aggregateID, ctx)
	if err != nil {
		return nil, err
	}

	return append(events, aggregate.LastEvent), nil
}
