package aggregate

import (
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/model"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"github.com/tangxusc/cqrs-db/pkg/snapshot"
	"strings"
	"sync/atomic"
	"time"
)

type Source struct {
	Id             string
	AggType        string
	Key            string
	locked         int32
	terminator     *time.Timer
	lifeTime       time.Duration
	Data           string
	LastUpdateTime time.Time
}

/*
停止的条件
1:10分钟未使用
*/
func (source *Source) start() {
	source.terminator = time.AfterFunc(source.lifeTime, source.gc)
}

func (source *Source) Locked() bool {
	if atomic.LoadInt32(&source.locked) != 0 {
		return true
	}
	return false
}

/*
1,根据聚合名称和id锁定聚合(lock),全局的聚合管理器
2,查询快照
3,根据快照查询事件
4,进行json merge溯源
5,返回结果
6,解锁(unlock)(事务提交)
*/
func (source *Source) Sourcing(handler *db.ConnHandler) (data string, err error) {
	handler.TxKey = source.Key
	source.lock()
	source.terminator.Reset(source.lifeTime)
	data, t := loadData(source)
	//根据快照的时间查询快照发生后的事件
	events, err := loadEvents(t, source.Id, source.AggType)
	if err != nil {
		return
	}
	//进行json merge溯源
	data, err = eventApply(data, events)
	if err != nil {
		return
	}
	go snapshot.Save(source.Id, source.AggType, data, events, time.Now())
	return
}

/*
加载顺序:
1,source对象
2,快照
*/
func loadData(source *Source) (string, time.Time) {
	var t = time.Time{}
	if !source.LastUpdateTime.IsZero() {
		t = source.LastUpdateTime
	}
	snap := loadSnapshot(source.Id, source.AggType, t)
	if snap != nil && !snap.CreateTime.IsZero() {
		source.Data = snap.Data
		source.LastUpdateTime = snap.CreateTime
		return snap.Data, snap.CreateTime
	}
	return source.Data, source.LastUpdateTime
}

func (source *Source) lock() {
	for {
		swapped := atomic.CompareAndSwapInt32(&source.locked, 0, 1)
		if swapped {
			return
		}
	}
}

func (source *Source) unlock() {
	for {
		swapped := atomic.CompareAndSwapInt32(&source.locked, 1, 0)
		if swapped {
			return
		}
	}
}

func (source *Source) gc() {
	if atomic.LoadInt32(&source.locked) != 0 {
		source.unlock()
	}
	SourceMap.Delete(source.Key)
	source.terminator.Stop()
}

func loadEvents(sh time.Time, id string, aggType string) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	newRow := func(types []*sql.ColumnType) []interface{} {
		e := &model.Event{}
		result := []interface{}{&e.Id, &e.Type, &e.AggId, &e.AggType, &e.CreateTime, &e.Data}
		events = append(events, e)
		return result
	}
	var err error
	if sh.IsZero() {
		err = proxy.QueryList(`select * from event where agg_id=? and agg_type=? order by create_time asc`, newRow, id, aggType)
	} else {
		err = proxy.QueryList(`select * from event where agg_id=? and agg_type=? and create_time > ? order by create_time asc`, newRow, id, aggType, sh)
	}
	if err != nil {
		return events, err
	}
	return events, nil
}

func loadSnapshot(id string, aggType string, t time.Time) *model.Snapshot {
	sh := &model.Snapshot{}
	var err error
	if t.IsZero() {
		err = proxy.QueryOne(`select * from snapshot where agg_id=? and agg_type=? order by create_time desc limit 0,1`, []interface{}{&sh.Id, &sh.AggId, &sh.AggType, &sh.CreateTime, &sh.Data}, id, aggType)
	} else {
		err = proxy.QueryOne(`select * from snapshot where agg_id=? and agg_type=? and create_time > ? order by create_time desc limit 0,1`, []interface{}{&sh.Id, &sh.AggId, &sh.AggType, &sh.CreateTime, &sh.Data}, id, aggType, t)
	}
	if err != nil {
		//没找到快照也进行聚合
		logrus.Warnf("[aggregate]查找快照出现错误,聚合:%v:%v,错误:%v", aggType, id, err)
	}
	return sh
}

func eventApply(sh string, events []*model.Event) (data string, err error) {
	data = sh
	if len(events) == 0 {
		return
	}
	//按照顺序合并数据
	for _, value := range events {
		data = data + value.Data
		data = strings.ReplaceAll(data, "}{", ",")
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(data), &result)
	bytes, err := json.Marshal(result)
	data = string(bytes)
	return
}
