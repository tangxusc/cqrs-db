package core

import (
	"time"
)

/*
快照
*/
type Snapshot struct {
	Id         string                 `bson:"_id"`
	AggId      string                 `bson:"agg_id"`
	AggType    string                 `bson:"agg_type"`
	CreateTime *time.Time             `bson:"create_time"`
	Data       map[string]interface{} `bson:"data"`
	Version    int                    `bson:"version"`
}
