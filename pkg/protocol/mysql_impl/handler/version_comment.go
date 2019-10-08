package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/parser"
	"github.com/xwb1989/sqlparser"
	"regexp"
)

func init() {
	variables := &versionComment{}
	compile, e := regexp.Compile(`(?i).*\s*select @@version_comment limit 1$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	mysql_impl.Handlers = append(mysql_impl.Handlers, variables)
}

type versionComment struct {
	compile *regexp.Regexp
}

func (s *versionComment) Match(stmt sqlparser.Statement) bool {
	sel, ok := stmt.(*sqlparser.Select)
	if ok {
		result := parser.ParseSelect(sel)
		columnLen := len(result.ColumnMap)
		for key, _ := range result.ColumnMap {
			if key == "@@version_comment" && columnLen == 1 {
				return true
			}
		}
	}
	return false
}

func (s *versionComment) Handler(query string, stmt sqlparser.Statement, handler *mysql_impl.ConnHandler) (*mysql.Result, error) {
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
