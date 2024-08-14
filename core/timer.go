package core

import "time"

type Timestamp = time.Time

type Timer interface {
	Now() Timestamp
}

type TimerUTC struct{}

func (t TimerUTC) Now() Timestamp {
	return time.Now().UTC()
}
