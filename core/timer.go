package core

import "time"

type Timer interface {
	Now() Timestamp
}

type DefaultTimer struct{}

func (t DefaultTimer) Now() Timestamp {
	return time.Now().UTC()
}
