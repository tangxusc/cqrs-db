package core

//SnapshotStore
//InitInterface

var eventRepository EventStore

/*
保存event
*/
type EventStore interface {
	SaveEvents(events Events) error
	UpdateEventStatus(event *Event, before MqStatus) error
	FindNotSentEventOrderByAsc() (Events, error)
}

func SetEventStore(r EventStore) {
	eventRepository = r
}

var aggregateCache AggregateCache

/*
具有淘汰机制,在内存不足时,淘汰某些aggregate
可根据key获取aggregate
*/
type AggregateCache interface {
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
