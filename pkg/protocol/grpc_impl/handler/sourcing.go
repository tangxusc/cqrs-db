package handler

import (
	"encoding/json"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/protocol/grpc_impl"
	"golang.org/x/net/context"
)

type SourcingHandler struct {
}

func NewSourcingHandler() *SourcingHandler {
	return &SourcingHandler{}
}

func (s *SourcingHandler) Sourcing(ctx context.Context, request *rpc.SourcingRequest) (response *rpc.SourcingResponse, e error) {
	data, version, e := core.Sourcing(request.AggId, request.AggType)
	if e != nil {
		return
	}
	response = &rpc.SourcingResponse{}
	response.AggId = request.AggId
	response.AggType = request.AggType
	response.Version = int32(version)
	bytes, e := json.Marshal(data)
	if e != nil {
		return
	}
	response.Data = string(bytes)
	return
}
