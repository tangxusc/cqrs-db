package handler

import (
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/mq"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl"
	"github.com/xwb1989/sqlparser"
	"testing"
	"time"
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
	result, e := variables.Handler(sqlString, statement, &mysql_impl.ConnHandler{})
	fmt.Println(result, e)
}

func TestFunc(t *testing.T) {
	event := core.NewEvent(``, ``, ``, ``, time.Now(), ``)
	retry := &Retry{}
	do := retry.Do(func(e error) {
		e = event.SuccessSend()
	})
	senderImpl := &mq.EventSenderImpl{}
	do = retry.Do(func(e error) {
		e = senderImpl.Send(nil)
	})
	fmt.Println(do)
}

type Retry struct {
	Count int
}

func (r *Retry) Do(f func(e error)) error {
	var e error
	f(e)
	if e == nil {
		return nil
	} else {
		r.Count = r.Count + 1
	}
	return e
}
