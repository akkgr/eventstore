package core

import "time"

type Timestamp = time.Time

type Timer interface {
	Now() Timestamp
}

type DefaultTimer struct{}

func (t DefaultTimer) Now() Timestamp {
	return time.Now().UTC()
}
