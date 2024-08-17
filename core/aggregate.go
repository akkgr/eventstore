package core

import "github.com/akkgr/eventstore/properties"

type Aggregate interface {
	Apply(event any) error
}

type AggregateBase struct {
	Id      string
	Entity  string
	Version int
	properties.Data
}

func (a *AggregateBase) New(id string, entity string) {
	a.Id = id
	a.Entity = entity
}

func (a *AggregateBase) SetVersion(version int) error {
	if version != a.Version+1 {
		return InvalidVersion{}
	}
	a.Version = version
	return nil
}
