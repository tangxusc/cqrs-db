package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/aggregate"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"regexp"
)

func init() {
	variables := &selectAggregate{}
	compile, e := regexp.Compile(`(?i).*\s*select \* from (\w+) where id='(\w+)'$`)
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

//TODO:实现聚合溯源
//TODO:实现sql查询别名的识别和应用
func (s *selectAggregate) Handler(query string) (*mysql.Result, error) {
	subMatch := s.compile.FindStringSubmatch(query)
	id, aggType, data, err := aggregate.Sourcing(subMatch)
	if err != nil {
		return nil, err
	}
	rows := make([][]interface{}, 0, 1)
	rows = append(rows, []interface{}{id, aggType, data})

	resultSet, err := mysql.BuildSimpleTextResultset([]string{"id", "type", "data"}, rows)
	if err != nil {
		return nil, err
	}
	return &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultSet,
	}, err
}
