package core

import (
	"encoding/json"
	"time"
)

type Event struct {
	Id      string    `json:"id"`
	Version int       `json:"version"`
	Entity  string    `json:"entity"`
	Action  string    `json:"action"`
	Created time.Time `json:"created"`
	Payload []byte    `json:"payload"`
}

func NewEvent(id string, version int, entity string, action string, payload any) (*Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Event{
		Id:      id,
		Version: version,
		Entity:  entity,
		Action:  action,
		Payload: data,
	}, nil
}
