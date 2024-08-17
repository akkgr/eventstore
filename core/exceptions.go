package core

type InvalidPayload struct {
}

func (e InvalidPayload) Error() string {
	return "Invalid payload"
}

type InvalidVersion struct {
}

func (e InvalidVersion) Error() string {
	return "Invalid version"
}

type EventsNotFound struct {
}

func (e EventsNotFound) Error() string {
	return "events not found"
}
