package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
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

func (s *insertEvent) Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	return event.Handler(query, stmt, handler)
}
