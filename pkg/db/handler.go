package db

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/sirupsen/logrus"
)

type QueryHandler interface {
	Match(query string) bool
	Handler(query string) (*mysql.Result, error)
}

var Handlers = make([]QueryHandler, 0)
var DefaultHandler QueryHandler

func NewHandler() server.Handler {
	return &TestHandler{}
}

type TestHandler struct {
}

func (t *TestHandler) UseDB(dbName string) error {
	logrus.Debug("UseDB.dbName:", dbName)
	return nil
}

func (t *TestHandler) HandleQuery(query string) (*mysql.Result, error) {
	logrus.Debug("HandleQuery.query:", query)
	for _, value := range Handlers {
		if value.Match(query) {
			return value.Handler(query)
		}
	}
	return DefaultHandler.Handler(query)
}

func (t *TestHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	logrus.Debugf("HandleFieldList(table:%v,fieldWildcard:%v)", table, fieldWildcard)
	return nil, fmt.Errorf("not supported now")
}

func (t *TestHandler) HandleStmtPrepare(query string) (params int, columns int, context interface{}, err error) {
	logrus.Debugf("HandleStmtPrepare(query:%v)", query)
	return 0, 0, nil, nil
}

//todo:无法返回结果,结果集为空(尽管已经返回了结果集)
func (t *TestHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	logrus.Debugf("HandleStmtExecute(context:%v,query:%v,args:%v)", context, query, args)
	return t.HandleQuery(query)
}

func (t *TestHandler) HandleStmtClose(context interface{}) error {
	logrus.Debugf("HandleStmtClose(context:%v)", context)
	return nil
}

func (t *TestHandler) HandleOtherCommand(cmd byte, data []byte) error {
	logrus.Debugf("HandleOtherCommand(cmd:%v,data:%v)", cmd, data)
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}
