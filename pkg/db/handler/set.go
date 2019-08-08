package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"github.com/xwb1989/sqlparser"
)

func init() {
	db.Handlers = append(db.Handlers, &set{})
}

type set struct {
}

func (s *set) Match(stmt sqlparser.Statement) bool {
	_, ok := stmt.(*sqlparser.Set)
	return ok
}

func (s *set) Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	_, _, err := proxy.Proxy(query)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
