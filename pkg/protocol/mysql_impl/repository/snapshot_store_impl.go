package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/util"
	"time"
)

type SnapshotStoreImpl struct {
	db *Conn
}

func NewSnapshotStoreImpl(db *Conn) *SnapshotStoreImpl {
	return &SnapshotStoreImpl{db: db}
}

func (s *SnapshotStoreImpl) FindLastOrderByCreateTimeDesc(aggId string, aggType string, start *time.Time) (*core.Snapshot, error) {
	sn := &core.Snapshot{}
	var err error
	if start == nil || start.IsZero() {
		err = s.db.QueryOne(`select * from snapshot where agg_id=? and agg_type=? order by create_time desc limit 0,1`,
			[]interface{}{&sn.Id, &sn.AggId, &sn.AggType, &sn.CreateTime, &sn.Data, &sn.Version}, aggId, aggType)
	} else {
		err = s.db.QueryOne(`select * from snapshot where agg_id=? and agg_type=? and create_time > ? order by create_time desc limit 0,1`,
			[]interface{}{&sn.Id, &sn.AggId, &sn.AggType, &sn.CreateTime, &sn.Data, &sn.Version}, aggId, aggType, start)
	}
	if err != nil {
		//没找到快照也进行聚合
		logrus.Warnf("[aggregate]查找快照出现错误,聚合:%v-%v,错误:%v", aggType, aggId, err)
	}
	return sn, nil
}

func (s *SnapshotStoreImpl) Save(aggId string, aggType string, cache *core.AggregateCache) {
	e := s.db.Exec(`insert into snapshot(id, agg_id, agg_type, create_time, data, revision) values (?, ?, ?, ?, ?,?)`,
		util.GenerateUuid(), aggId, aggType, cache.UpdateTime, cache.Data, cache.Version)
	if e != nil {
		logrus.Warnf("[snapshot]保存快照失败,聚合:%v-%v,错误:%v", aggType, aggId, e)
	}
}
