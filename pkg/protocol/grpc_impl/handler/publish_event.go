package handler

import (
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/protocol/grpc_impl"
	"github.com/tangxusc/cqrs-db/pkg/util"
	"golang.org/x/net/context"
	"time"
)

type EventPublishHandler struct {
}

func NewEventPublishHandler() *EventPublishHandler {
	return &EventPublishHandler{}
}

func (h *EventPublishHandler) Publish(ctx context.Context, request *rpc.PublishRequest) (response *rpc.PublishResponse, e error) {
	events := buildEvents(request)
	e = events.SaveAndSend()
	if e != nil {
		return
	}
	response = &rpc.PublishResponse{}
	response.Version = request.Version
	return
}

func buildEvents(request *rpc.PublishRequest) core.Events {
	events := make([]*core.Event, 0)
	event := core.NewEvent(util.GenerateUuid(), request.EventType, request.AggId, request.AggType, time.Now(), request.Data, int(request.Version))
	events = append(events, event)
	return events
}
