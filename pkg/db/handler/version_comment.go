package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"regexp"
)

func init() {
	variables := &versionComment{}
	compile, e := regexp.Compile(`(?i).*\s*select @@version_comment limit 1$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

type versionComment struct {
	compile *regexp.Regexp
}

func (s *versionComment) Match(query string) bool {
	return s.compile.MatchString(query)
}

func (s *versionComment) Handler(query string) (*mysql.Result, error) {
	//mysql> select @@version_comment limit 1;
	//	+------------------------------+
	//	| @@version_comment            |
	//		+------------------------------+
	//	| MySQL Community Server - GPL |
	//		+------------------------------+
	//		1 row in set (0.00 sec)

	var resultset *mysql.Resultset
	var err error
	rows := make([][]interface{}, 0, 1)
	rows = append(rows, []interface{}{"MySQL Community Server - GPL"})

	resultset, err = mysql.BuildSimpleTextResultset([]string{"@@version_comment"}, rows)

	result := &mysql.Result{
		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
		InsertId:     0,
		AffectedRows: 0,
		Resultset:    resultset,
	}

	return result, err
}
