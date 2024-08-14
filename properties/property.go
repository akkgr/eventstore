package properties

import (
	"time"

	"github.com/akkgr/eventstore/core"
	"github.com/shopspring/decimal"
)

type Value = string
type ValueId = string
type ValueType string

const (
	Text   ValueType = "text"
	Number ValueType = "number"
	Date   ValueType = "date"
)

type Property interface {
	GetId() ValueId
	GetType() ValueType
	GetValue() (interface{}, error)
}

type property struct {
	Id    ValueId   `json:"id" dynamodbav:"id"`
	Type  ValueType `json:"type" dynamodbav:"type"`
	Value Value     `json:"value" dynamodbav:"value"`
}

func (p *property) GetId() ValueId {
	return p.Id
}

func (p *property) GetType() ValueType {
	return p.Type
}

func (p *property) GetValue() (interface{}, error) {
	switch p.Type {
	case Text:
		return p.Value, nil
	case Number:
		return decimal.NewFromString(p.Value)
	case Date:
		return time.Parse(time.RFC3339, p.Value)
	default:
		return nil, nil
	}
}

func NewTextProperty(id ValueId, v core.Text) Property {
	return &property{
		Id:    id,
		Type:  Text,
		Value: v,
	}
}

func NewNumberProperty(id ValueId, v core.Number) Property {
	return &property{
		Id:    id,
		Type:  Text,
		Value: v.String(),
	}
}

func NewDateProperty(id ValueId, v core.Timestamp) Property {
	return &property{
		Id:    id,
		Type:  Text,
		Value: v.Format(time.RFC3339),
	}
}

type Properties = map[string]Property
type Collection = map[string]Properties
type Collections = map[string]Collection
type Data struct {
	Properties  Properties  `json:"properties" dynamodbav:"properties"`
	Collections Collections `json:"collections" dynamodbav:"collections"`
}
