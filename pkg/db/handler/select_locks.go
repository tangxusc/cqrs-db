package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/aggregate"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/db/parser"
	"github.com/xwb1989/sqlparser"
	"strings"
)

type selectLocks struct {
}

func (s *selectLocks) Match(stmt sqlparser.Statement) bool {
	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return false
	}
	result := parser.ParseSelect(sel)
	//table名称为 xxx_aggregate格式
	return strings.ToLower(result.TableName) == "locks_agg"
}

func (s *selectLocks) Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	var resultset *mysql.Resultset
	var err error
	rows := make([][]interface{}, 0)
	aggregate.SourceMap.Range(func(key, value interface{}) bool {
		entry := value.(*aggregate.Source)
		if entry.Locked() {
			rows = append(rows, []interface{}{entry.Key, "locked"})
		}
		return true
	})

	resultset, err = mysql.BuildSimpleTextResultset([]string{"key", "lock"}, rows)

	result := &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultset,
	}

	return result, err
}

func init() {
	db.Handlers = append(db.Handlers, &selectLocks{})
}
