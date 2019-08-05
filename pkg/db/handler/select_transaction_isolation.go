package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"regexp"
)

func init() {
	variables := &transactionIsolation{}
	compile, e := regexp.Compile(`(?i).*\s*SELECT @@session.transaction_isolation$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

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
	var resultset *mysql.Resultset
	var err error
	rows := make([][]interface{}, 0, 1)
	rows = append(rows, []interface{}{"REPEATABLE-READ"})

	resultset, err = mysql.BuildSimpleTextResultset([]string{"@@session.transaction_isolation"}, rows)

	result := &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultset,
	}

	return result, err
}
