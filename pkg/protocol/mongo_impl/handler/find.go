package handler

import (
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/core"
	protocol "github.com/tangxusc/mongo-protocol"
	"strings"
)

type Handler struct {
}

func (q *Handler) Process(header *protocol.MsgHeader, r *protocol.Reader, conn *protocol.ConnContext) error {
	query := &protocol.Query{}
	if e := query.UnMarshal(r); e != nil {
		return e
	}
	id, aggType, e := getAggInfo(query)
	if e != nil {
		return e
	}
	data, version, e := core.Sourcing(id, aggType)
	if e != nil {
		return e
	}
	reply := protocol.NewReply(header.RequestID)
	reply.NumberReturned = 1
	data["version"] = version
	reply.Documents = data
	e = reply.Write(conn)
	if e != nil {
		return e
	}
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
