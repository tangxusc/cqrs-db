package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/proxy"
	"github.com/xwb1989/sqlparser"
)

func init() {
	mysql_impl.Handlers = append(mysql_impl.Handlers, &set{})
}

type set struct {
}

func (s *set) Match(stmt sqlparser.Statement) bool {
	_, ok := stmt.(*sqlparser.Set)
	return ok
}

func (s *set) Handler(query string, stmt sqlparser.Statement, handler *mysql_impl.ConnHandler) (*mysql.Result, error) {
	_, _, err := proxy.Proxy(query)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
