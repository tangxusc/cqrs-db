package core

import (
	"time"
)

type MqStatus string

const (
	NotSend MqStatus = "NotSend"
	Sent    MqStatus = "Sent"
)

/*
事件
*/
type Event struct {
	Id         string    `bson:"id"`
	Type       string    `bson:"event_type"`
	AggId      string    `bson:"agg_id"`
	AggType    string    `bson:"agg_type"`
	CreateTime time.Time `bson:"create_time"`
	Data       string    `bson:"data"`
	Status     MqStatus  `bson:"status"`
	Version    int       `bson:"version"`
}

func (event *Event) SuccessSend() error {
	event.Status = Sent
	return eventRepository.UpdateEventStatus(event, NotSend)
}

func NewEvent(id string, eventType string, aggId string, aggType string, createTime time.Time, data string, version int) *Event {
	return &Event{Id: id, Type: eventType, AggId: aggId, AggType: aggType, CreateTime: createTime, Data: data, Status: NotSend, Version: version}
}

type Events []*Event

func (events Events) SaveAndSend() error {
	//获取aggregate,并放入聚合发送队列,再由聚合发送出
	event := events[0]
	agg := aggregateCache.Get(event.AggId, event.AggType)

	return agg.PutSendChan(events)
}

type KeyBuild func(e *Event) string

func (events Events) Group(keyBuild KeyBuild) map[string]Events {
	result := make(map[string]Events)
	for _, v := range events {
		key := keyBuild(v)
		value, ok := result[key]
		if !ok {
			value = make([]*Event, 0)
		}
		value = append(value, v)
		result[key] = value
	}
	return result
}

func (events Events) SendToRecovery() error {
	event := events[0]
	agg := aggregateCache.Get(event.AggId, event.AggType)
	return agg.PutRecoveryChan(events)
}

func (events Events) ToEventArray() []*Event {
	return events
}

//TODO:可以删除
func (events Events) Length() int {
	return len([]*Event(events))
}
