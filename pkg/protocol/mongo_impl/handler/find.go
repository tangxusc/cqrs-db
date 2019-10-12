package handler

import (
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/core"
	protocol "github.com/tangxusc/mongo-protocol"
	"strings"
)

type FindHandler struct {
}

func (q *FindHandler) Support(query *protocol.Query) bool {
	fullCollectionName := strings.ToLower(query.FullCollectionName)
	contains := strings.Contains(fullCollectionName, "_aggregate")
	return contains
}

func (q *FindHandler) Process(query *protocol.Query, reply *protocol.Reply) error {
	id, aggType, e := getAggInfo(query)
	if e != nil {
		return e
	}
	data, version, e := core.Sourcing(id, aggType)
	if e != nil {
		return e
	}
	reply.NumberReturned = 1
	data["version"] = version
	reply.Documents = data
	return nil
}

func getAggInfo(query *protocol.Query) (string, string, error) {
	id, ok := query.Query["id"]
	if !ok {
		return "", "", fmt.Errorf("[mongodb]缺少查询id")
	}
	fullCollectionName := strings.ToLower(query.FullCollectionName)
	contains := strings.Contains(fullCollectionName, "_aggregate")
	if !contains {
		return "", "", fmt.Errorf("[mongodb]聚合名称错误")
	}
	collectionName := strings.ReplaceAll(fullCollectionName, "_aggregate", "")
	return id.(string), collectionName, nil
}

func NewFindHandler() *FindHandler {
	return &FindHandler{}
}
