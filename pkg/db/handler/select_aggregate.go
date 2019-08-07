package handler

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/aggregate"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/db/parser"
	"regexp"
)

var ColumnsName = []string{"id", "agg_type", "data"}

func init() {
	variables := &selectAggregate{}
	compile, e := regexp.Compile(`(?i).*\s*select\s+.+\s+from (\w+)_Aggregate\s*(?:\w+)*\s*where id='(\w+)'$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

type selectAggregate struct {
	compile *regexp.Regexp
}

func (s *selectAggregate) Match(query string) bool {
	return s.compile.MatchString(query)
}

//example: select id as c1 , agg_type as c2, data as c3 from test_aggregate a1 where id='1'
func (s *selectAggregate) Handler(query string) (*mysql.Result, error) {
	result, e := parser.ParseSql(query)
	if e != nil {
		return nil, e
	}
	parseResult := result.(*parser.SelectParseResult)
	for key, _ := range parseResult.ColumnMap {
		if key != ColumnsName[0] && key != ColumnsName[1] && key != ColumnsName[2] {
			return nil, fmt.Errorf("不支持的列名称:%v", key)
		}
	}

	id, aggType, e := ParseIdAndType(s, query)
	if e != nil {
		return nil, e
	}
	data, e := aggregate.Sourcing(id, aggType)
	if e != nil {
		return nil, e
	}
	rows := make([][]interface{}, 0, 1)
	rows = append(rows, []interface{}{id, aggType, data})

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
	strings := make([]string, 3)
	for key, _ := range strings {
		strings[key] = result.ColumnMap[ColumnsName[key]]
	}
	return strings
}

func ParseIdAndType(s *selectAggregate, query string) (id string, aggType string, err error) {
	subMatch := s.compile.FindStringSubmatch(query)
	if len(subMatch) != 3 {
		err = fmt.Errorf("sql语句错误,请传入正确的参数")
		return
	}
	id = subMatch[2]
	aggType = subMatch[1]
	return
}
