package model

import "time"

type Snapshot struct {
	Id         string
	AggId      string
	AggType    string
	CreateTime time.Time
	Data       string
}
