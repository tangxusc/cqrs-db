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
	Id         string
	Type       string
	AggId      string
	AggType    string
	CreateTime time.Time
	Data       string
	Status     MqStatus
}

func (event *Event) SuccessSend() error {
	event.Status = Sent
	return eventRepository.UpdateEventStatus(event, NotSend)
}

func NewEvent(id string, eventType string, aggId string, aggType string, createTime time.Time, data string) *Event {
	return &Event{Id: id, Type: eventType, AggId: aggId, AggType: aggType, CreateTime: createTime, Data: data, Status: NotSend}
}

type Events []*Event

func (events Events) SaveAndSend() error {
	//保存event
	err := eventRepository.SaveEvents(events)
	if err != nil {
		return err
	}
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
			value = make([]*Event, 1)
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
