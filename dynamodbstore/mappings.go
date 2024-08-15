package dynamodbstore

import "github.com/akkgr/eventstore/eventstore"

func eventFrom(e eventstore.Event) event {
	return event{
		Id:      e.Id,
		Version: e.Version,
		Entity:  e.Entity,
		Action:  e.Action,
		Created: e.Created,
		Data:    e.Data,
	}
}

func eventTo(e event) eventstore.Event {
	return eventstore.Event{
		Id:      e.Id,
		Version: e.Version,
		Entity:  e.Entity,
		Action:  e.Action,
		Created: e.Created,
		Data:    e.Data,
	}
}

func aggregateFrom(a eventstore.Aggregate) aggregate {
	return aggregate{
		Id:        a.Id,
		Entity:    a.Entity,
		Created:   a.Created,
		LastEvent: eventFrom(a.LastEvent),
	}
}

func aggregateTo(a aggregate) eventstore.Aggregate {
	return eventstore.Aggregate{
		Id:        a.Id,
		Entity:    a.Entity,
		Created:   a.Created,
		LastEvent: eventTo(a.LastEvent),
	}
}

func snapshotFrom(s eventstore.Snapshot) snapshot {
	return snapshot{
		Id:      s.Id,
		Version: s.Version,
		Entity:  s.Entity,
		Created: s.Created,
		Data:    s.Data,
	}
}

func snapshotTo(s snapshot) eventstore.Snapshot {
	return eventstore.Snapshot{
		Id:      s.Id,
		Version: s.Version,
		Entity:  s.Entity,
		Created: s.Created,
		Data:    s.Data,
	}
}
