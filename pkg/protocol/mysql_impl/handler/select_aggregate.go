package handler

import (
	"encoding/json"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/parser"
	"github.com/xwb1989/sqlparser"
	"strings"
)

var ColumnsName = []string{"id", "agg_type", "data"}

func init() {
	mysql_impl.Handlers = append(mysql_impl.Handlers, &selectAggregate{})
}

type selectAggregate struct {
}

func (s *selectAggregate) Match(stmt sqlparser.Statement) bool {
	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return false
	}
	result := parser.ParseSelect(sel)
	if len(result.Where) != 1 {
		return false
	}
	//table名称为 xxx_aggregate格式
	return strings.HasSuffix(strings.ToLower(result.TableName), "_aggregate")
}

//example: select id as c1 , agg_type as c2, data as c3 from test_aggregate a1 where id='1'
//example: select * from test_aggregate a1 where id='1'
//example: select * from test_aggregate where id='1'
func (s *selectAggregate) Handler(query string, stmt sqlparser.Statement, handler *mysql_impl.ConnHandler) (*mysql.Result, error) {
	parseResult := parser.ParseSelect(stmt.(*sqlparser.Select))
	for key := range parseResult.ColumnMap {
		if key != ColumnsName[0] && key != ColumnsName[1] && key != ColumnsName[2] {
			return nil, fmt.Errorf("不支持的列名称:%v", key)
		}
	}
	id, aggType := ParseIdAndType(parseResult)

	data, version, e := core.Sourcing(id, aggType)
	data["version"] = version
	bytes, e := json.Marshal(data)
	if e != nil {
		return nil, e
	}

	rows := make([][]interface{}, 0, 1)
	rows = append(rows, []interface{}{id, aggType, string(bytes)})

	resultSet, e := mysql.BuildSimpleTextResultset(getColumn(parseResult), rows)
	if e != nil {
		return nil, e
	}
	return &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultSet,
	}, e
}

func getColumn(result *parser.SelectParseResult) []string {
	if len(result.ColumnMap) == 0 {
		return ColumnsName
	}
	returns := make([]string, len(ColumnsName))
	for key := range returns {
		returns[key] = result.ColumnMap[ColumnsName[key]]
	}
	return returns
}

func ParseIdAndType(result *parser.SelectParseResult) (id string, aggType string) {
	id = result.Where[0]
	aggType = strings.ReplaceAll(result.TableName, "_aggregate", "")
	aggType = strings.ReplaceAll(aggType, "_Aggregate", "")
	return
}
