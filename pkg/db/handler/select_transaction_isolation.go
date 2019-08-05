package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"regexp"
)

//func init() {
//	variables := &transactionIsolation{}
//	compile, e := regexp.Compile(`(?i).*\s*SELECT @@session.transaction_isolation$`)
//	if e != nil {
//		panic(e.Error())
//	}
//	variables.compile = compile
//	db.Handlers = append(db.Handlers, variables)
//}

type transactionIsolation struct {
	compile *regexp.Regexp
}

func (s *transactionIsolation) Match(query string) bool {
	return s.compile.MatchString(query)
}

func (s *transactionIsolation) Handler(query string) (*mysql.Result, error) {
	//mysql> SELECT @@session.transaction_isolation;
	//+---------------------------------+
	//| @@session.transaction_isolation |
	//	+---------------------------------+
	//| REPEATABLE-READ                 |
	//	+---------------------------------+
	//	1 row in set (0.00 sec)

	columnNames, columnValues, e := proxy.Proxy(query)
	if e != nil {
		return nil, e
	}
	resultset, e := mysql.BuildSimpleTextResultset(columnNames, columnValues)
	if e != nil {
		return nil, e
	}

	return &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultset,
	}, e
}
