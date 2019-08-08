package aggregate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/model"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"github.com/tangxusc/cqrs-db/pkg/snapshot"
	"strings"
	"time"
)

/*
1,根据聚合名称和id锁定聚合(lock),全局的聚合管理器
2,查询快照
3,根据快照查询事件
4,进行json merge溯源
5,返回结果
6,解锁(unlock)
*/
//TODO:加入缓存
func Sourcing(id string, aggType string, handler *db.ConnHandler) (data string, err error) {
	if !handler.TxBegin {
		err = fmt.Errorf("需要开启事务才能查询聚合")
		return
	}
	if len(handler.TxKey) != 0 {
		err = fmt.Errorf("一个事务中只能查询一次聚合")
		return
	}
	//lock
	handler.TxKey = Lock(id, aggType)
	//在事务结束时,解锁
	//TODO:超时,30秒后事务任然未解锁,则主动解锁
	//查询快照
	sh := loadSnapshot(id, aggType)
	//根据快照的时间查询快照发生后的事件
	events, err := loadEvents(sh, id, aggType)
	if err != nil {
		return
	}
	//进行json merge溯源
	data, err = eventApply(sh, events)
	if err != nil {
		return
	}
	go snapshot.Save(id, aggType, data, events, time.Now())
	return
}

func loadEvents(sh *model.Snapshot, id string, aggType string) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	newRow := func(types []*sql.ColumnType) []interface{} {
		e := &model.Event{}
		result := []interface{}{&e.Id, &e.Type, &e.AggId, &e.AggType, &e.CreateTime, &e.Data}
		events = append(events, e)
		return result
	}
	var err error
	if sh.CreateTime.IsZero() {
		err = proxy.QueryList(`select * from event where agg_id=? and agg_type=? order by create_time asc`, newRow, id, aggType)
	} else {
		err = proxy.QueryList(`select * from event where agg_id=? and agg_type=? and create_time > ? order by create_time asc`, newRow, id, aggType, sh.CreateTime)
	}
	if err != nil {
		return events, err
	}
	return events, nil
}

func loadSnapshot(id string, aggType string) *model.Snapshot {
	sh := &model.Snapshot{}
	err := proxy.QueryOne(`select * from snapshot where agg_id=? and agg_type=? order by create_time desc limit 0,1`, []interface{}{&sh.Id, &sh.AggId, &sh.AggType, &sh.CreateTime, &sh.Data}, id, aggType)
	if err != nil {
		//没找到快照也进行聚合
		logrus.Warnf("[aggregate]查找快照出现错误,聚合:%v:%v,错误:%v", aggType, id, err)
	}
	return sh
}

func eventApply(sh *model.Snapshot, events []*model.Event) (data string, err error) {
	//快照不存在
	if sh.CreateTime.IsZero() {
		data = ""
	} else {
		data = sh.Data
	}
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
