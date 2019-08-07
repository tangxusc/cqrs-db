package parser

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

type ParseResult interface {
	ToString() string
}

/*
解析sql
在此只解析select和insert语句
*/
func ParseSql(sql string) (result ParseResult, err error) {
	statement, err := sqlparser.Parse(sql)
	if err != nil {
		return
	}
	switch stmt := statement.(type) {
	case *sqlparser.Select:
		result = parseSelect(stmt)
	case *sqlparser.Insert:
		//TODO:暂未实现
		fmt.Println(stmt)
	default:
		err = fmt.Errorf("只支持select,insert语句")
	}
	return
}
