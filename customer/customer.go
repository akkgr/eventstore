package customer

import (
	"encoding/json"

	"github.com/akkgr/eventstore/core"
	"github.com/akkgr/eventstore/properties"
	"github.com/google/uuid"
)

type Customer struct {
	core.AggregateBase
	Name   string
	Status string
}

const (
	CustomerCreated = "CustomerCreated"
	CustomerUpdated = "CustomerUpdated"
)

type CustomerCreatedEvent struct {
	Name string
}

type CustomerUpdatedEvent struct {
	Name   string
	Status string
	properties.Data
}

func (c *Customer) customerCreated(e core.Event) error {
	var payload CustomerCreatedEvent
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return err
	}
	c.New(uuid.New().String(), "Customer")
	c.Name = payload.Name

	return nil
}

func (c *Customer) customerUpdated(e core.Event) error {
	var payload CustomerUpdatedEvent
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return err
	}
	c.Name = payload.Name
	c.Status = payload.Status
	c.Data = payload.Data

	return nil
}

func (c *Customer) Apply(e core.Event) error {

	if err := c.SetVersion(e.Version); err != nil {
		return err
	}

	switch e.Action {
	case CustomerCreated:
		return c.customerCreated(e)
	case CustomerUpdated:
		return c.customerUpdated(e)
	default:
		return core.InvalidPayload{}
	}
}
