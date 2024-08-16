package properties

type InvalidPropertyType struct {
}

func (e InvalidPropertyType) Error() string {
	return "Invalid property type"
}
