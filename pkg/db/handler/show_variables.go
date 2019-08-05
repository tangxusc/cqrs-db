package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	py "github.com/tangxusc/cqrs-db/pkg/proxy"
	"regexp"
)

func init() {
	variables := &showVariables{}
	compile, e := regexp.Compile(`(?i).*\s*SHOW VARIABLES$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

type showVariables struct {
	compile *regexp.Regexp
}

func (s *showVariables) Match(query string) bool {
	return s.compile.MatchString(query)
}

func (s *showVariables) Handler(query string) (*mysql.Result, error) {
	var resultSet *mysql.Resultset
	var err error

	var name string
	var value string
	rows := make([][]interface{}, 0)
	err = py.Query(query, func() []interface{} {
		return []interface{}{&name, &value}
	}, func(row []interface{}) {
		i := make([]interface{}, 0, len(row))
		i = append(i, name)
		i = append(i, value)
		rows = append(rows, i)
	})
	if err != nil {
		return nil, err
	}

	resultSet, err = mysql.BuildSimpleTextResultset([]string{"Variable_name", "Value"}, rows)

	result := &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultSet,
	}

	return result, err
}
