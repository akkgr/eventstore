package core

import "time"

type Timer interface {
	Now() time.Time
}

type DefaultTimer struct{}

func (t DefaultTimer) Now() time.Time {
	return time.Now().UTC()
}
