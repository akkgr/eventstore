package eventstore

type InvalidVersion struct {
}

func (e InvalidVersion) Error() string {
	return "version mismatch"
}
