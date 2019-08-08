package handler

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/aggregate"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/xwb1989/sqlparser"
)

func init() {
	db.Handlers = append(db.Handlers, &transaction{})
}

type transaction struct {
}

func (s *transaction) Match(stmt sqlparser.Statement) bool {
	_, ok := stmt.(*sqlparser.Begin)
	if ok {
		return true
	}
	_, ok = stmt.(*sqlparser.Commit)
	if ok {
		return true
	}
	_, ok = stmt.(*sqlparser.Rollback)
	if ok {
		return true
	}

	return false
}

func (s *transaction) Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	switch stmt.(type) {
	case *sqlparser.Begin:
		if handler.TxBegin {
			return nil, fmt.Errorf("已经开启了事务")
		}
		handler.TxBegin = true
	case *sqlparser.Commit, *sqlparser.Rollback:
		if !handler.TxBegin {
			return nil, fmt.Errorf("未开启事务,无法提交事务")
		}
		handler.TxBegin = false
		aggregate.UnLockWithKey(handler.TxKey)
		handler.TxKey = ""
	}
	return nil, nil
}
