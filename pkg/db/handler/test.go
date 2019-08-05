package handler

import (
	"cqrs-db/pkg/db"
	"github.com/siddontang/go-mysql/mysql"
	"regexp"
	"strings"
)

//TODO:删除测试对象
func init() {
	variables := &testQuery{}
	db.Handlers = append(db.Handlers, variables)
}

type testQuery struct {
	compile *regexp.Regexp
}

func (s *testQuery) Match(query string) bool {
	return strings.Contains(query, "some_table_name")
}

func (s *testQuery) Handler(query string) (*mysql.Result, error) {
	var resultset *mysql.Resultset
	var err error
	row := make([][]interface{}, 0, 1)
	row = append(row, []interface{}{`{"id": 1, "name": "lnmp.cn"}`, "2"})
	resultset, err = mysql.BuildSimpleTextResultset([]string{"test2_0_", "test1_0_"}, row)
	result := &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultset,
	}

	return result, err
}
