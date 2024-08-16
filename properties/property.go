package properties

import (
	"time"

	"github.com/shopspring/decimal"
)

type Property struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PropertyTypes interface {
	string | decimal.Decimal | time.Time
}

func NewProperty[T PropertyTypes](id string, v T) (Property, error) {
	val := any(v)
	switch val := val.(type) {
	case string:
		return Property{
			Id:    id,
			Type:  "text",
			Value: val,
		}, nil
	case decimal.Decimal:
		return Property{
			Id:    id,
			Type:  "number",
			Value: val.String(),
		}, nil
	case time.Time:
		return Property{
			Id:    id,
			Type:  "date",
			Value: val.Format(time.RFC3339),
		}, nil
	default:
		return Property{}, InvalidPropertyType{}
	}
}

func (p *Property) GetValue() (any, error) {
	switch p.Type {
	case "text":
		return p.Value, nil
	case "number":
		return decimal.NewFromString(p.Value)
	case "date":
		return time.Parse(time.RFC3339, p.Value)
	default:
		return nil, InvalidPropertyType{}
	}
}

type Properties = map[string]Property
type Collection = map[string]Properties
type Collections = map[string]Collection
type Data struct {
	Properties  Properties  `json:"properties"`
	Collections Collections `json:"collections"`
}
