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

type MongoAggregate struct {
	Id      string        `bson:"_id"`
	AggType string        `bson:"agg_type"`
	Events  []*core.Event `bson:"events"`
}

/*
保存聚合事件
*/
func (s *EventStoreImpl) SaveEvents(agg *core.Aggregate, events core.Events) error {
	var update = true
	andUpdate := s.Database(s.dbName).Collection(s.collectionName).FindOneAndUpdate(s.ctx,
		bson.M{"_id": agg.AggId, "agg_type": agg.AggType},
		bson.M{"$push": bson.M{"events": bson.M{"$each": events.ToEventArray()}}},
		&options.FindOneAndUpdateOptions{
			Upsert: &update,
		})
	return andUpdate.Err()
}

func (s *EventStoreImpl) UpdateEventStatus(event *core.Event, before core.MqStatus) error {
	result, e := s.Database(s.dbName).Collection(s.collectionName).UpdateOne(s.ctx,
		bson.M{
			"_id":      event.AggId,
			"agg_type": event.AggType,
			"events": bson.M{
				"$elemMatch": bson.M{"id": event.Id, "status": before},
			},
		},
		bson.M{
			"$set": bson.M{
				"events.$.status": event.Status,
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
			"events": bson.M{
				"$elemMatch": bson.M{"status": core.NotSend},
			},
		},
	)
	//TODO:查询投影options.FindOptions{}
	if e != nil {
		return nil, e
	}
	events := make([]*core.Event, 0)
	for cursor.Next(s.ctx) {
		agg := &MongoAggregate{}
		e := cursor.Decode(agg)
		if e != nil {
			return nil, e
		}
		for _, event := range agg.Events {
			if event.Status == core.NotSend {
				events = append(events, event)
			}
		}
	}

	return events, nil
}

func (s *EventStoreImpl) FindByOrderByAsc(aggId string, aggType string, time *time.Time) (core.Events, error) {
	//{$match: {user_id : A_id} },
	//{$unwind:"events"},
	//{$match: {'bonus.type' : b} },
	//cursor, e := s.Database(s.dbName).Collection(s.collectionName).Aggregate(s.ctx, )
	panic("implement me")
}
