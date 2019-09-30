package event

import (
	"encoding/json"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/db/parser"
	"github.com/tangxusc/cqrs-db/pkg/util"
	"github.com/xwb1989/sqlparser"
	"time"
)

type Impl struct {
	ImplId         string    `json:"id"`
	ImplEventType  string    `json:"event_type"`
	ImplAggId      string    `json:"agg_id"`
	ImplAggType    string    `json:"agg_type"`
	ImplCreateTime time.Time `json:"create_time"`
	ImplData       string    `json:"data"`
}

func (impl *Impl) Data() string {
	return impl.ImplData
}

func (impl *Impl) ToJson() ([]byte, error) {
	return json.Marshal(impl)
}

func (impl *Impl) CreateTime() time.Time {
	return impl.ImplCreateTime
}

func (impl *Impl) AggType() string {
	return impl.ImplAggType
}

func (impl *Impl) AggId() string {
	return impl.ImplAggId
}

func (impl *Impl) EventType() string {
	return impl.ImplEventType
}

func (i *Impl) Id() string {
	return i.ImplId
}

func NewEventImpl(id string, eventType string, aggId string, aggType string, createTime time.Time, data string) *Impl {
	return &Impl{ImplId: id, ImplEventType: eventType, ImplAggId: aggId, ImplAggType: aggType, ImplCreateTime: createTime, ImplData: data}
}

type Event interface {
	Id() string
	EventType() string
	AggId() string
	AggType() string
	CreateTime() time.Time
	Data() string
	ToJson() ([]byte, error)
}

func BuildEvents(parseResult *parser.InsertParseResult) (events []Event) {
	events = make([]Event, len(parseResult.Values))
	for k := range parseResult.Values {
		id := util.GenerateUuid()
		eventType := string(parseResult.ValueMaps["type"][k].([]byte))
		aggId := string(parseResult.ValueMaps["agg_id"][k].([]byte))
		aggType := string(parseResult.ValueMaps["agg_type"][k].([]byte))
		createTimeString := string(parseResult.ValueMaps["create_time"][k].([]byte))
		createTime, e := time.Parse(`2006-01-02 15:04:05`, createTimeString)
		fmt.Println(createTime, e)
		data := string(parseResult.ValueMaps["data"][k].([]byte))
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
	events := BuildEvents(result)
	//3,保存event
	err = SaveEvent(events)
	if err != nil {
		return nil, err
	}
	//4,发送到mq
	err = SenderImpl.SendEvents(events)
	if err != nil {
		return nil, err
	}

	return nil, err
}

var SenderImpl Sender

type Sender interface {
	SendEvents([]Event) error
}
