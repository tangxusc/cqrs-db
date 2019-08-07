package model

import "time"

type Event struct {
	Id         string
	Type       string
	AggId      string
	AggType    string
	CreateTime time.Time
	Data       string
}
