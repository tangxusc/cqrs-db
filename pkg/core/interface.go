package core

import (
	"time"
)

//InitInterface
//DbRetryer.NewInstance()
//DbRetryer.Do(f){
//if e!=nil{
//	f()
//}
//}

var eventRepository EventStore

/*
保存event
*/
type EventStore interface {
	SaveEvents(aggregate *Aggregate, events Events) error
	UpdateEventStatus(event *Event, before MqStatus) error
	FindNotSentEventOrderByAsc() (Events, error)
	FindByOrderByAsc(aggId string, aggType string, time *time.Time) (Events, error)
}

func SetEventStore(r EventStore) {
	eventRepository = r
}

var aggregateCache AggregateManager

func SetAggregateManager(a AggregateManager) {
	aggregateCache = a
}

/*
具有淘汰机制,在内存不足时,淘汰某些aggregate
可根据key获取aggregate
map 定长(先入先出) 时间过期
*/
type AggregateManager interface {
	Get(aggId, aggType string) *Aggregate
}

var eventSender EventSender

/*
发送消息
*/
type EventSender interface {
	Send(event *Event) error
}

func SetEventSender(r EventSender) {
	eventSender = r
}

type SnapshotStore interface {
	//按照创建时间降序查找最后一个快照
	FindLastOrderByCreateTimeDesc(aggId string, aggType string, start *time.Time) (*Snapshot, error)
	Save(aggId string, aggType string, cache *AggregateCache)
}

var snapshotStore SnapshotStore

func SetSnapshotStore(s SnapshotStore) {
	snapshotStore = s
}

type SnapshotSaveStrategy interface {
	Allow(aggId string, aggType string, data map[string]interface{}, events Events) bool
}

func SetSnapshotSaveStrategy(s SnapshotSaveStrategy) {
	snapshotSaveStrategy = s
}

var snapshotSaveStrategy SnapshotSaveStrategy
