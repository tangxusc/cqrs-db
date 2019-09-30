package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Restorer struct {
	ticket *time.Ticker
}

/*
   TODO:recovery 执行周期时间
*/
func NewRestorer() *Restorer {
	return &Restorer{}
}

func (r *Restorer) Start(ctx context.Context) {
	recovery()
	r.ticket = time.NewTicker(time.Second * 5)
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
