package core

import "time"

/*
快照
*/
type Snapshot struct {
	Id         string
	AggId      string
	AggType    string
	CreateTime time.Time
	Data       string
}
