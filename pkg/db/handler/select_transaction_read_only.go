package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"regexp"
)

func init() {
	variables := &transactionReadOnly{}
	compile, e := regexp.Compile(`(?i).*\s*SELECT @@session.transaction_read_only$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

type transactionReadOnly struct {
	compile *regexp.Regexp
}

func (s *transactionReadOnly) Match(query string) bool {
	return s.compile.MatchString(query)
}

func (s *transactionReadOnly) Handler(query string) (*mysql.Result, error) {
	//mysql> SELECT @@session.transaction_read_only;
	//+---------------------------------+
	//| @@session.transaction_read_only |
	//	+---------------------------------+
	//|                               0 |
	//	+---------------------------------+
	//	1 row in set (0.07 sec)
	var resultset *mysql.Resultset
	var err error
	rows := make([][]interface{}, 0, 1)
	rows = append(rows, []interface{}{0})

	resultset, err = mysql.BuildSimpleTextResultset([]string{"@@session.transaction_read_only"}, rows)

	result := &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultset,
	}

	return result, err
}
