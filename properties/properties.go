package properties

import (
	"time"

	"github.com/shopspring/decimal"
)

type Property struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewTextProperty(v string) Property {
	return Property{
		Id:    "",
		Type:  "text",
		Value: v,
	}
}

func NewNumberProperty(v decimal.Decimal) Property {
	return Property{
		Id:    "",
		Type:  "number",
		Value: v.String(),
	}
}

func NewDateProperty(v time.Time) Property {
	return Property{
		Id:    "",
		Type:  "date",
		Value: v.Format(time.RFC3339),
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
