package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"time"
)

type Restorer struct {
	ticket *time.Ticker
}

func NewRestorer() *Restorer {
	return &Restorer{}
}

func (r *Restorer) Start(ctx context.Context) {
	if eventSender == nil {
		logrus.Warnf("[recovery] EventSender 未配置,将不发送事件")
		return
	}
	recovery()
	r.ticket = time.NewTicker(time.Second * time.Duration(config.Instance.ServerDb.RecoveryInterval))
	go func() {
		for {
			select {
			case <-ctx.Done():
				r.ticket.Stop()
				return
			case <-r.ticket.C:
				recovery()
			}
		}
	}()
}

func recovery() {
	events, e := eventRepository.FindNotSentEventOrderByAsc()
	if e != nil {
		logrus.Errorf("[recovery] FindNotSentEventOrderByAsc error:%v", e)
	}
	groupEvent := events.Group(key)
	for _, v := range groupEvent {
		send(v)
	}
}

func key(e *Event) string {
	return e.AggId + `-` + e.AggType
}

func send(v Events) {
	go func() {
		e := v.SendToRecovery()
		if e != nil {
			logrus.Errorf("[recovery] SendToRecovery error:%v", e)
		}
	}()
}
