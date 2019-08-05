package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"regexp"
)

func init() {
	variables := &selectAggregate{}
	compile, e := regexp.Compile(`(?i).*\s*select * from (\w+) where id=(\w+)$`)
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

//todo:实现聚合溯源
func (s *selectAggregate) Handler(query string) (*mysql.Result, error) {
	//var name interface{}
	//var value string
	//var err error
	//rows := make([][]interface{}, 0, 1)
	//
	//err = py.Query(query, func() []interface{} {
	//	return []interface{}{&name, &value}
	//}, func(row []interface{}) {
	//	i := make([]interface{}, 0, len(row))
	//	i = append(i, name)
	//	i = append(i, value)
	//	rows = append(rows, i)
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//resultSet, err := mysql.BuildSimpleTextResultset([]string{"id", "type", "data"}, rows)
	//
	//return &mysql.Result{
	//	Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
	//	InsertId:     0,
	//	AffectedRows: 0,
	//	Resultset:    resultSet,
	//}, err
	return nil, nil
}
