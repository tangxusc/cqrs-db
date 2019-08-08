package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"github.com/xwb1989/sqlparser"
)

type defaultHandler struct {
}

func (d *defaultHandler) Match(stmt sqlparser.Statement) bool {
	return true
}

func (d *defaultHandler) Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	columnNames, columnValues, err := proxy.Proxy(query)
	if err != nil {
		return nil, err
	}
	resultSet, err := mysql.BuildSimpleTextResultset(columnNames, columnValues)
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

func init() {
	db.DefaultHandler = &defaultHandler{}
}
