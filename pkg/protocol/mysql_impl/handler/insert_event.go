package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/parser"
	"github.com/tangxusc/cqrs-db/pkg/util"
	"github.com/xwb1989/sqlparser"
	"strings"
	"time"
)

var Columns = []string{"id", "type", "agg_id", "agg_type", "create_time", "data"}

func init() {
	mysql.Handlers = append(mysql.Handlers, &insertEvent{})
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

func (s *insertEvent) Handler(query string, stmt sqlparser.Statement, handler *mysql_impl.ConnHandler) (*mysql.Result, error) {
	result, e := parser.ParseInsert(stmt.(*sqlparser.Insert))
	if e != nil {
		return nil, e
	}
	events, e := BuildEvents(result)
	if e != nil {
		return nil, e
	}
	return nil, events.SaveAndSend()
}

func BuildEvents(parseResult *parser.InsertParseResult) (events core.Events, e error) {
	events = make([]*core.Event, len(parseResult.Values))
	for k := range parseResult.Values {
		id := util.GenerateUuid()
		eventType := string(parseResult.ValueMaps["type"][k].([]byte))
		aggId := string(parseResult.ValueMaps["agg_id"][k].([]byte))
		aggType := string(parseResult.ValueMaps["agg_type"][k].([]byte))
		createTimeString := string(parseResult.ValueMaps["create_time"][k].([]byte))
		var createTime time.Time
		createTime, e = time.Parse(`2006-01-02 15:04:05`, createTimeString)
		if e != nil {
			return
		}
		data := string(parseResult.ValueMaps["data"][k].([]byte))
		events[k] = core.NewEvent(id, eventType, aggId, aggType, createTime, data)
	}
	return events, nil
}
