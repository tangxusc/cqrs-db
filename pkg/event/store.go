package event

import (
	"database/sql"
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
)

var Columns = []string{"id", "type", "agg_id", "agg_type", "create_time", "data", "mq_status"}

func SaveEvent(events []Event) error {
	if len(events) == 0 {
		return fmt.Errorf("events不能为空")
	}
	return proxy.Tx(func(tx *sql.Tx) error {
		for _, v := range events {
			stmt, e := tx.Prepare(`insert into event(id, type, agg_id, agg_type, create_time, data, mq_status) values (?,?,?,?,?,?,?)`)
			if e != nil {
				return e
			}
			_, e = stmt.Exec(v.Id(), v.EventType(), v.AggId(), v.AggType(), v.CreateTime(), v.Data(), NotSend)
			if e != nil {
				return e
			}
			return stmt.Close()
		}
		return nil
	})
}
