package db

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/sirupsen/logrus"
	"github.com/xwb1989/sqlparser"
)

type QueryHandler interface {
	Match(stmt sqlparser.Statement) bool
	Handler(query string, stmt sqlparser.Statement, handler *ConnHandler) (*mysql.Result, error)
}

var Handlers = make([]QueryHandler, 0)
var DefaultHandler QueryHandler

func NewHandler() server.Handler {
	return &ConnHandler{}
}

type ConnHandler struct {
	TxBegin bool
	TxKey   string
}

func (t *ConnHandler) UseDB(dbName string) error {
	logrus.Debug("UseDB.dbName:", dbName)
	return nil
}

func (t *ConnHandler) HandleQuery(query string) (*mysql.Result, error) {
	logrus.Debug("HandleQuery.query:", query)
	statement, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}
	for _, value := range Handlers {
		if value.Match(statement) {
			return value.Handler(query, statement, t)
		}
	}
	return DefaultHandler.Handler(query, statement, t)
}

func (t *ConnHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	logrus.Debugf("HandleFieldList(table:%v,fieldWildcard:%v)", table, fieldWildcard)
	return nil, fmt.Errorf("not supported now")
}

func (t *ConnHandler) HandleStmtPrepare(query string) (params int, columns int, context interface{}, err error) {
	logrus.Debugf("HandleStmtPrepare(query:%v)", query)
	return 0, 0, nil, nil
}

//todo:无法返回结果,结果集为空(尽管已经返回了结果集)
func (t *ConnHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	logrus.Debugf("HandleStmtExecute(context:%v,query:%v,args:%v)", context, query, args)
	return t.HandleQuery(query)
}

func (t *ConnHandler) HandleStmtClose(context interface{}) error {
	logrus.Debugf("HandleStmtClose(context:%v)", context)
	return nil
}

func (t *ConnHandler) HandleOtherCommand(cmd byte, data []byte) error {
	logrus.Debugf("HandleOtherCommand(cmd:%v,data:%v)", cmd, data)
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}
