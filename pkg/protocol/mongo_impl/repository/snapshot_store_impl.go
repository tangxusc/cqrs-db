package repository

import (
	"context"
	"github.com/siddontang/go/bson"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SnapshotStoreImpl struct {
	*MongoConn
	ctx            context.Context
	collectionName string
}

func NewSnapshotStoreImpl(mongoConn *MongoConn, ctx context.Context, collectionName string) *SnapshotStoreImpl {
	return &SnapshotStoreImpl{MongoConn: mongoConn, ctx: ctx, collectionName: collectionName}
}

func (s *SnapshotStoreImpl) FindLastOrderByCreateTimeDesc(aggId string, aggType string, start *time.Time) (*core.Snapshot, error) {
	var filter bson.M
	if start == nil || start.IsZero() {
		filter = bson.M{
			"agg_id":   aggId,
			"agg_type": aggType,
		}
	} else {
		filter = bson.M{
			"agg_id":   aggId,
			"agg_type": aggType,
			"create_time": bson.M{
				"$gt": start,
			},
		}
	}
	one := s.Database(s.dbName).Collection(s.collectionName).FindOne(s.ctx,
		filter,
		&options.FindOneOptions{
			Sort: bson.M{
				"create_time": 1,
			},
		},
	)
	snapshot := &core.Snapshot{}
	if e := one.Decode(snapshot); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, e
	}
	return snapshot, nil
}

func (s *SnapshotStoreImpl) Save(aggId string, aggType string, cache *core.AggregateCache) {
	snapshot := &core.Snapshot{
		Id:         util.GenerateUuid(),
		AggId:      aggId,
		AggType:    aggType,
		CreateTime: cache.UpdateTime,
		Data:       cache.Data,
		Version:    cache.Version,
	}
	result, e := s.Database(s.dbName).Collection(s.collectionName).InsertOne(s.ctx, snapshot)
	if e != nil {
		logrus.Errorf("[snapshot] save snapshot[%s-%s] error:%v ", aggType, aggId, e)
	} else {
		logrus.Debugf("[snapshot] save snapshot success :%s", result.InsertedID)
	}

}
