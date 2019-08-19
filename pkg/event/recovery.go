package event

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"time"
)

var ticket *time.Timer

func RecoveryEvent(ctx context.Context) {
	ticket = time.AfterFunc(time.Second*5, recoveryEvent)
}

/*
从数据库中获取记录,并调用sender发送
*/
func recoveryEvent() {
	events := make([]Event, 0)
	newRow := func(types []*sql.ColumnType) []interface{} {
		e := &Impl{}
		result := []interface{}{&e.ImplId, &e.ImplEventType, &e.ImplAggId, &e.ImplAggType, &e.ImplCreateTime, &e.ImplData}
		events = append(events, e)
		return result
	}
	e := proxy.QueryList(`select id,type,agg_id,agg_type,create_time,data from event where mq_status=? order by create_time asc`, newRow, NotSend)
	if e != nil {
		logrus.Warnf("[event]获取未发送event出现错误,%v", e)
		return
	}

	e = SenderImpl.SendEvents(events)
	if e != nil {
		logrus.Warnf("[event]重新发送event出现错误,%v", e)
	}
}

func Stop() {
	ticket.Stop()
}
