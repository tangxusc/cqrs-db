package proxy

import (
	"context"
	"cqrs-db/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"os"
)

var db *sql.DB

func InitConn(ctx context.Context) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=true", config.Instance.Proxy.Username, config.Instance.Proxy.Password,
		"tcp", config.Instance.Proxy.Address, config.Instance.Proxy.Port, config.Instance.Proxy.Database)
	var e error
	db, e = sql.Open("mysql", dsn)
	if e != nil {
		logrus.Errorf("[proxy]连接出现错误,url:%v,错误:%v", dsn, e.Error())
		os.Exit(1)
	}
	defer db.Close()

	<-ctx.Done()
}

func Query(query string, newRow func() []interface{}, rowAfter func(row []interface{})) error {
	logrus.Debugf("[proxy]查询:%s", query)
	rows, e := db.Query(query)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		row := newRow()
		e := rows.Scan(row...)
		if e != nil {
			return e
		}
		rowAfter(row)
	}
	return nil
}
