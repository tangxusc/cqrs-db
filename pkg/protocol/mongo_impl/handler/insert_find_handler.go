package handler

import (
	"fmt"
	protocol "github.com/tangxusc/mongo-protocol"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type InsertFindHandler struct {
	insert *InsertHandler
}

func NewInsertFindHandler(insert *InsertHandler) *InsertFindHandler {
	return &InsertFindHandler{insert: insert}
}

func (q *InsertFindHandler) Support(query *protocol.Query) bool {
	fullCollectionName := strings.ToLower(query.FullCollectionName)
	if fullCollectionName == `aggregate.$cmd` {
		_, ok := query.Query[`insert`]
		if ok {
			return true
		}
	}
	return false
}

//db.a1_aggregate.find({'id':'4'})
func (q *InsertFindHandler) Process(query *protocol.Query, reply *protocol.Reply) error {
	v := query.Query["insert"]
	docs := query.Query["documents"]
	i := docs.([]interface{})
	ms := make([]bson.M, len(i))
	for k, v := range i {
		ms[k] = v.(bson.M)
		fmt.Println(ms, k, v)
	}
	insert := &protocol.Insert{
		FullCollectionName: v.(string),
		Documents:          ms,
	}
	e := q.insert.HandlerInsert(insert)
	if e != nil {
		reply.ResponseFlags = protocol.QueryFailure
		reply.NumberReturned = 1
		reply.Documents = map[string]interface{}{"$err": fmt.Sprintf("%v", e)}
		return nil
	}
	reply.NumberReturned = 1
	reply.Documents = map[string]interface{}{"ok": 1}
	return nil
}
