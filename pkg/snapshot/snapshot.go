package snapshot

import (
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/model"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"github.com/tangxusc/cqrs-db/pkg/util"
	"time"
)

func Save(aggId, aggType string, data string, events []*model.Event, now time.Time) {
	if Strategy.Allow(aggId, aggType, data, events) {
		SaveSnapshot(aggId, aggType, data, now)
	}
}

func SaveSnapshot(aggId string, aggType string, data string, time time.Time) {
	e := proxy.Exec(`insert into snapshot(id, agg_id, agg_type, create_time, data) values (?, ?, ?, ?, ?)`, util.GenerateUuid(), aggId, aggType, time, data)
	if e != nil {
		logrus.Warnf("[snapshot]保存快照失败,聚合:%v:%v,错误:%v", aggType, aggId, e)
	}
}

//默认提供基于事件数量的快照
var Strategy SaveStrategy = &CountStrategy{Max: 100}

type SaveStrategy interface {
	Allow(Id string, aggType string, data string, events []*model.Event) bool
}

type CountStrategy struct {
	Max int
}

func (c *CountStrategy) Allow(Id string, aggType string, data string, events []*model.Event) bool {
	if len(events) > c.Max {
		return true
	}
	return false
}
