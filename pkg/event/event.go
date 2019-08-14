package event

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/db/parser"
	"github.com/xwb1989/sqlparser"
	"time"
)

type Impl struct {
	id         string
	eventType  string
	aggId      string
	aggType    string
	createTime time.Time
	data       string
}

func (impl *Impl) Data() string {
	return impl.data
}

func (impl *Impl) CreateTime() time.Time {
	return impl.createTime
}

func (impl *Impl) AggType() string {
	return impl.aggType
}

func (impl *Impl) AggId() string {
	return impl.aggId
}

func (impl *Impl) EventType() string {
	return impl.eventType
}

func (i *Impl) Id() string {
	return i.id
}

func NewEventImpl(id string, eventType string, aggId string, aggType string, createTime time.Time, data string) *Impl {
	return &Impl{id: id, eventType: eventType, aggId: aggId, aggType: aggType, createTime: createTime, data: data}
}

type Event interface {
	Id() string
	EventType() string
	AggId() string
	AggType() string
	CreateTime() time.Time
	Data() string
}

func BuildEvents(parseResult *parser.InsertParseResult) (events []Event) {
	events = make([]Event, len(parseResult.ValueMaps))
	for k, _ := range parseResult.Values {
		id := parseResult.ValueMaps["id"][k].(string)
		eventType := parseResult.ValueMaps["eventType"][k].(string)
		aggId := parseResult.ValueMaps["aggId"][k].(string)
		aggType := parseResult.ValueMaps["aggType"][k].(string)
		createTime := parseResult.ValueMaps["createTime"][k].(time.Time)
		data := parseResult.ValueMaps["data"][k].(string)
		events[k] = NewEventImpl(id, eventType, aggId, aggType, createTime, data)
	}
	return events
}

/*
1,检查sql
2,保存event到数据库
3,go 发送到mq
3,go grpc推送event
4,返回保存成功
*/
func Handler(query string, stmt sqlparser.Statement, handler *db.ConnHandler) (*mysql.Result, error) {
	var err error
	//1,必须在一个事务中
	if !handler.TxBegin {
		return nil, fmt.Errorf("必须在事务中才能发布事件")
	}
	//2,分析出数据
	result, err := parser.ParseInsert(stmt.(*sqlparser.Insert))
	if err != nil {
		return nil, err
	}
	//3,保存event
	err = SaveEvent(result)
	if err != nil {
		return nil, err
	}
	events := BuildEvents(result)
	//4,发送到mq
	err = SenderImpl.SendEvents(events)
	if err != nil {
		return nil, err
	}
	//TODO:5,leader grpc 推送event

	return nil, err
}

var SenderImpl Sender

type Sender interface {
	SendEvents([]Event) error
}
