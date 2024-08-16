package eventstore

import (
	"time"
)

type Event struct {
	Id      string    `json:"id"`
	Version int       `json:"version"`
	Entity  string    `json:"entity"`
	Action  string    `json:"action"`
	Created time.Time `json:"created"`
	Data    []byte    `json:"data"`
}
