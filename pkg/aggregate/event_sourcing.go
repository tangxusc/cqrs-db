package aggregate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"strings"
	"time"
)

type snapshot struct {
	Id         string
	AggId      string
	AggType    string
	createTime time.Time
	data       string
}
type event struct {
	Id         string
	AggId      string
	AggType    string
	createTime time.Time
	data       string
}

/*
1,分析出聚合名称和id
根据聚合名称和id锁定聚合(lock),全局的聚合管理器
2,查询快照
3,根据快照查询事件
4,进行json merge溯源
5,返回结果
解锁(unlock)
*/
//TODO:加入缓存
func Sourcing(ident []string) (id string, aggType string, data string, err error) {
	//检查是否传入了id和类型
	if len(ident) != 2 {
		err = fmt.Errorf("sql语句错误,请传入正确的参数")
		return
	}
	id = ident[1]
	aggType = ident[0]
	//lock
	Lock(id, aggType)
	defer UnLock(id, aggType)
	//查询快照
	sh := &snapshot{}
	err = proxy.QueryOne(`select * from snapshot where agg_id=? and agg_type=? order by create_time limit 0,1`, []interface{}{&sh.Id, &sh.AggId, &sh.AggType, &sh.createTime, &sh.data}, id, aggType)
	if err != nil {
		return
	}
	//根据快照的时间查询快照发生后的事件
	events := make([]*event, 0)
	err = proxy.QueryList(`select * from event where agg_id=? and agg_type=? order by create_time asc`, func(types []*sql.ColumnType) []interface{} {
		e := &event{}
		result := []interface{}{&e.Id, &e.AggId, &e.AggType, &e.createTime, &e.data}
		events = append(events, e)
		return result
	})
	if err != nil {
		return
	}
	//进行json merge溯源
	data, err = eventApply(sh, events)
	//TODO:发送到快照,根据策略创建新的快照
	return
}

func eventApply(sh *snapshot, events []*event) (data string, err error) {
	//快照不存在
	if sh.createTime.IsZero() {
		data = ""
	} else {
		data = sh.data
	}
	//按照顺序合并数据
	for _, value := range events {
		data = data + value.data
		data = strings.ReplaceAll(data, "}{", ",")
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(data), &result)
	bytes, err := json.Marshal(result)
	data = string(bytes)
	return
}
