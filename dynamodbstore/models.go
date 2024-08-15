package dynamodbstore

import "time"

type event struct {
	Id      string    `json:"id" dynamodbav:"id"`
	Version int       `json:"version" dynamodbav:"version"`
	Entity  string    `json:"entity" dynamodbav:"entity"`
	Action  string    `json:"action" dynamodbav:"action"`
	Created time.Time `json:"created" dynamodbav:"created"`
	Data    []byte    `json:"data" dynamodbav:"data"`
}

type aggregate struct {
	Id        string    `json:"id" dynamodbav:"id"`
	Entity    string    `json:"entity" dynamodbav:"entity"`
	Created   time.Time `json:"created" dynamodbav:"created"`
	LastEvent event     `json:"lastEvent" dynamodbav:"lastEvent"`
}

type snapshot struct {
	Id      string    `json:"id" dynamodbav:"id"`
	Version int       `json:"varsion" dynamodbav:"version"`
	Entity  string    `json:"entity" dynamodbav:"entity"`
	Created time.Time `json:"created" dynamodbav:"created"`
	Data    []byte    `json:"data" dynamodbav:"data"`
}
