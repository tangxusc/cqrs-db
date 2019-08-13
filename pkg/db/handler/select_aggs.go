package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/aggregate"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/db/parser"
	"github.com/xwb1989/sqlparser"
	"strings"
)

type selectAggs struct {
}

func (s *selectAggs) Match(stmt sqlparser.Statement) bool {
	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return false
	}
	result := parser.ParseSelect(sel)
	//table名称为 agg_info
	return strings.ToLower(result.TableName) == "agg_info"
}

func (s *selectAggs) Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	var resultset *mysql.Resultset
	var err error
	rows := make([][]interface{}, 0)
	aggregate.SourceMap.Range(func(key, value interface{}) bool {
		entry := value.(*aggregate.Source)
		var lockstatus string
		if entry.Locked() {
			lockstatus = "locked"
		} else {
			lockstatus = "unlock"
		}
		rows = append(rows, []interface{}{entry.Key, lockstatus, entry.Data, entry.LastUpdateTime.String()})
		return true
	})

	resultset, err = mysql.BuildSimpleTextResultset([]string{"key", "locked", "data", "last_update_time"}, rows)

	result := &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultset,
	}

	return result, err
}

func init() {
	db.Handlers = append(db.Handlers, &selectAggs{})
}
