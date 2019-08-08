package handler

import (
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/xwb1989/sqlparser"
	"testing"
)

func TestHandler(t *testing.T) {
	variables := &selectAggregate{}

	sqlString := `select * from User_Aggregate where id='xudslajsdfhsadhf'`
	statement, e := sqlparser.Parse(sqlString)
	if e != nil {
		panic(e.Error())
	}
	match := variables.Match(statement)
	if !match {
		panic("不匹配语句")
	}
	result, e := variables.Handler(sqlString, statement, &db.ConnHandler{})
	fmt.Println(result, e)
}
