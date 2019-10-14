package handler

import (
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/util"
	protocol "github.com/tangxusc/mongo-protocol"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"time"
)

type InsertHandler struct {
	compile *regexp.Regexp
}

func (i *InsertHandler) Process(header *protocol.MsgHeader, r *protocol.Reader, conn *protocol.ConnContext) error {
	insert := &protocol.Insert{}
	e := insert.UnMarshal(r)
	if e != nil {
		return e
	}
	e = i.validateDocuments(insert)
	if e != nil {
		return e
	}
	events, e := i.buildEvents(insert)
	if e != nil {
		return e
	}
	e = events.SaveAndSend()
	if e != nil {
		return e
	}
	return nil
}

func (i *InsertHandler) validateDocuments(insert *protocol.Insert) error {
	match := i.compile.MatchString(insert.FullCollectionName)
	if !match {
		return fmt.Errorf("[mongo]collection格式错误,事件collection必须以'_event'结尾")
	}
	submatch := i.compile.FindStringSubmatch(insert.FullCollectionName)
	if len(submatch) != 2 {
		return fmt.Errorf("[mongo]collection格式错误,事件collection必须以'_event'结尾")
	}
	for k, doc := range insert.Documents {
		e := i.validateDocument(doc, k)
		if e != nil {
			return e
		}
	}
	return nil
}

//aggId,eventType,data,version,createTime
func (i *InsertHandler) validateDocument(doc bson.M, k int) error {
	_, ok := doc["aggId"]
	if !ok {
		return fmt.Errorf("[mongo]第%v文档中,缺少字段aggId", k)
	}
	_, ok = doc["eventType"]
	if !ok {
		return fmt.Errorf("[mongo]第%v文档中,缺少字段eventType", k)
	}
	_, ok = doc["data"]
	if !ok {
		return fmt.Errorf("[mongo]第%v文档中,缺少字段data", k)
	}
	_, ok = doc["version"]
	if !ok {
		return fmt.Errorf("[mongo]第%v文档中,缺少字段version", k)
	}
	_, ok = doc["createTime"]
	if !ok {
		return fmt.Errorf("[mongo]第%v文档中,缺少字段createTime", k)
	}

	return nil
}

func (i *InsertHandler) buildEvents(insert *protocol.Insert) (core.Events, error) {
	events := make([]*core.Event, len(insert.Documents))
	aggType := i.compile.FindStringSubmatch(insert.FullCollectionName)[1]
	for k, doc := range insert.Documents {
		val := doc["createTime"].(string)
		createTime, e := time.Parse(`2006-01-02 15:04:05`, val)
		if e != nil {
			return nil, e
		}
		fmt.Println(doc["eventType"])
		bytes, e := bson.MarshalJSON(doc["data"])
		if e != nil {
			return nil, e
		}
		event := core.NewEvent(util.GenerateUuid(), doc["eventType"].(string), doc["aggId"].(string), aggType, createTime, string(bytes), int(doc["version"].(float64)))
		events[k] = event
	}
	return events, nil
}

func NewInsertHandler() *InsertHandler {
	compile, e := regexp.Compile(`(?i)(.+)_event$`)
	if e != nil {
		panic(fmt.Errorf(`[mongo]protocol格式检查正则错误,%w`, e))
	}
	return &InsertHandler{
		compile: compile,
	}
}
