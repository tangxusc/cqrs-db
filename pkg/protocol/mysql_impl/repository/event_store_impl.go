package repository

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"time"
)

type EventStoreImpl struct {
	db *Conn
}

func (s *EventStoreImpl) FindByOrderByAsc(aggId string, aggType string, t *time.Time) (core.Events, error) {
	var err error
	events := make([]*core.Event, 0)
	newRow := func(types []*sql.ColumnType) []interface{} {
		e := &core.Event{}
		result := []interface{}{&e.Id, &e.Type, &e.AggId, &e.AggType, &e.CreateTime, &e.Data}
		events = append(events, e)
		return result
	}
	if t == nil || t.IsZero() {
		err = s.db.QueryList(`select * from event where agg_id=? and agg_type=? order by create_time asc`, newRow, aggId, aggType)
	} else {
		err = s.db.QueryList(`select * from event where agg_id=? and agg_type=? and create_time > ? order by create_time asc`, newRow, aggId, aggType, t)
	}
	if err != nil {
		return events, err
	}
	return events, nil
}

func (s *EventStoreImpl) FindNotSentEventOrderByAsc() (core.Events, error) {
	events := make([]*core.Event, 0)
	newRow := func(types []*sql.ColumnType) []interface{} {
		e := &core.Event{}
		result := []interface{}{&e.Id, &e.Type, &e.AggId, &e.AggType, &e.CreateTime, &e.Data, &e.Status}
		events = append(events, e)
		return result
	}
	e := s.db.QueryList(`select id,type,agg_id,agg_type,create_time,data,mq_status from event where mq_status=? order by create_time asc`, newRow, core.NotSend)
	if e != nil {
		logrus.Warnf("[EventStoreImpl]获取未发送event出现错误,%v", e)
		return nil, e
	}
	return events, nil
}

func NewEventStoreImpl(c *Conn) *EventStoreImpl {
	return &EventStoreImpl{c}
}

func (s *EventStoreImpl) UpdateEventStatus(event *core.Event, before core.MqStatus) error {
	e := s.db.Exec(`update event set mq_status=? where agg_id=? and agg_type=? and mq_status=?`, event.Status, event.AggId, event.AggType, before)
	if e != nil {
		logrus.Errorf("[EventStoreImpl]事件已发送,更新数据库发生错误:%s,%s,错误:%v", event.AggId, event.AggType, e)
		return e
	}
	return nil
}

func (s *EventStoreImpl) SaveEvents(events core.Events) error {
	if len(events) == 0 {
		return fmt.Errorf("events不能为空")
	}
	data := make([][]interface{}, len(events))
	for k, v := range events {
		data[k] = []interface{}{v.Id, v.Type, v.AggId, v.AggType, v.CreateTime, v.Data, core.NotSend}
	}
	return s.db.Execs(`into event(id, type, agg_id, agg_type, create_time, data, mq_status) values (?,?,?,?,?,?,?)`, data)
}
