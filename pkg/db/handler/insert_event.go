package handler

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/db/parser"
	"github.com/tangxusc/cqrs-db/pkg/event"
	"github.com/xwb1989/sqlparser"
	"strings"
)

var Columns = []string{"id", "type", "agg_id", "agg_type", "create_time", "data"}

func init() {
	db.Handlers = append(db.Handlers, &insertEvent{})
}

type insertEvent struct {
}

//insert into event(type, agg_id, agg_type, create_time, data)
//values ('E1', '1', 'A1', '', '{"name":"test1"}')
func (s *insertEvent) Match(stmt sqlparser.Statement) bool {
	insert, ok := stmt.(*sqlparser.Insert)
	if !ok {
		return false
	}
	if strings.ToLower(insert.Table.Name.String()) != "event_aggregate" {
		return false
	}

	return true
}

/*
1,检查sql
2,保存event到数据库
3,go 发送到mq
3,go grpc推送event
4,返回保存成功
*/
func (s *insertEvent) Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	var err error
	//1,必须在一个事务中
	if !handler.TxBegin {
		return nil, fmt.Errorf("必须在事务中才能发布事件")
	}
	//2,分析出数据
	result, err := parser.ParseInsert(stmt.(*sqlparser.Insert))
	if err != nil {
		return nil, err
	}
	//3,保存event
	err = event.SaveEvent(result)
	if err != nil {
		return nil, err
	}
	//TODO:4,发送到mq
	//TODO:5,leader grpc 推送event

	return nil, err
}
