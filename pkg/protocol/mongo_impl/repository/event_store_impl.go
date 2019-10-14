package repository

import (
	"context"
	"github.com/siddontang/go/bson"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type EventStoreImpl struct {
	*MongoConn
	ctx            context.Context
	collectionName string
}

func NewEventStoreImpl(mongoConn *MongoConn, ctx context.Context, collectionName string) *EventStoreImpl {
	return &EventStoreImpl{MongoConn: mongoConn, ctx: ctx, collectionName: collectionName}
}

/*
保存聚合事件
*/
func (s *EventStoreImpl) SaveEvents(agg *core.Aggregate, events core.Events) error {
	es := make([]interface{}, len(events))
	for key, value := range events {
		es[key] = value
	}
	_, e := s.Database(s.dbName).Collection(s.collectionName).InsertMany(s.ctx, es)
	if e != nil {
		return e
	}
	return nil
}

func (s *EventStoreImpl) UpdateEventStatus(event *core.Event, before core.MqStatus) error {
	result, e := s.Database(s.dbName).Collection(s.collectionName).UpdateOne(s.ctx,
		bson.M{
			"_id": event.Id,
		},
		bson.M{
			"$set": bson.M{
				"status": event.Status,
			},
		})
	if e != nil {
		return e
	}
	logrus.Debugf("[mongodb]updated event count:%v", result.ModifiedCount)
	return nil
}

func (s *EventStoreImpl) FindNotSentEventOrderByAsc() (core.Events, error) {
	cursor, e := s.Database(s.dbName).Collection(s.collectionName).Find(s.ctx,
		bson.M{
			"status": core.NotSend,
		},
		&options.FindOptions{
			Sort: bson.M{
				"create_time": 1,
			},
		},
	)
	if e != nil {
		return nil, e
	}
	events := make([]*core.Event, 0)
	for cursor.Next(s.ctx) {
		event := &core.Event{}
		e := cursor.Decode(event)
		if e != nil {
			return nil, e
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *EventStoreImpl) FindByOrderByAsc(aggId string, aggType string, time *time.Time) (core.Events, error) {
	var filter bson.M
	if time == nil || time.IsZero() {
		filter = bson.M{
			"agg_id":   aggId,
			"agg_type": aggType,
		}
	} else {
		filter = bson.M{
			"agg_id":   aggId,
			"agg_type": aggType,
			"create_time": bson.M{
				"$gt": time,
			},
		}
	}
	cursor, e := s.Database(s.dbName).Collection(s.collectionName).Find(s.ctx,
		filter,
		&options.FindOptions{
			Sort: bson.M{
				"create_time": 1,
			},
		},
	)
	if e != nil {
		return nil, e
	}
	events := make([]*core.Event, 0)
	for cursor.Next(s.ctx) {
		event := &core.Event{}
		e := cursor.Decode(event)
		if e != nil {
			return nil, e
		}
		events = append(events, event)
	}
	return events, nil
}
