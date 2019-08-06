package proxy

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"os"
	"time"
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
	db.SetConnMaxLifetime(time.Duration(config.Instance.Proxy.LifeTime) * time.Second)
	db.SetMaxOpenConns(config.Instance.Proxy.MaxOpen)
	db.SetMaxIdleConns(config.Instance.Proxy.MaxIdle)
	<-ctx.Done()
}

func QueryOne(sqlString string, scan []interface{}, param ...interface{}) error {
	stmt, e := db.Prepare(sqlString)
	if e != nil {
		return e
	}
	defer stmt.Close()
	row := stmt.QueryRow(param...)
	e = row.Scan(scan...)
	//未找到记录
	if e != nil && e == sql.ErrNoRows {
		return nil
	}
	if e != nil {
		return e
	}
	return nil
}

func QueryList(sqlString string, newRow func(types []*sql.ColumnType) []interface{}) error {
	return Query(sqlString, newRow, func(row []interface{}) {
		//忽略
	}, func(strings []string) {
		//忽略
	})
}

func Query(query string, newRow func(types []*sql.ColumnType) []interface{}, rowAfter func(row []interface{}), setColumnNames func([]string)) error {
	logrus.Debugf("[proxy]Query:%s", query)
	rows, e := db.Query(query)
	if e != nil {
		return e
	}
	defer rows.Close()
	types, e := rows.ColumnTypes()
	if e != nil {
		return e
	}
	strings, e := rows.Columns()
	if e != nil {
		return e
	}
	setColumnNames(strings)
	for rows.Next() {
		row := newRow(types)
		e := rows.Scan(row...)
		if e != nil {
			return e
		}
		rowAfter(row)
	}
	return nil
}
func Proxy(query string) (columnNames []string, columnValues [][]interface{}, err error) {
	logrus.Debugf("[proxy]Proxy:%s", query)
	var temp interface{} = ""
	var rowOrigin []interface{}
	var result []interface{}

	columnValues = make([][]interface{}, 0)
	err = Query(query,
		func(types []*sql.ColumnType) []interface{} {
			if result == nil {
				result = make([]interface{}, len(types))
				rowOrigin = make([]interface{}, 0, len(types))
				for key, _ := range types {
					rowOrigin = append(rowOrigin, temp)
					result[key] = &rowOrigin[key]
				}
			}
			return result
		},
		func(row []interface{}) {
			i := make([]interface{}, len(row))
			for key, _ := range row {
				i[key] = rowOrigin[key]
			}
			columnValues = append(columnValues, i)
		},
		func(strings []string) {
			columnNames = strings
		})
	if err != nil {
		return
	}
	return
}
