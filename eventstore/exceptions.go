package eventstore

type InvalidVersion struct {
}

func (e InvalidVersion) Error() string {
	return "version mismatch"
}

type EventsNotFound struct {
}

func (e EventsNotFound) Error() string {
	return "events not found"
}
