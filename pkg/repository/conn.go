package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"time"
)

type Conn struct {
	*sql.DB
}

func InitConn(ctx context.Context) (*Conn, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=true", config.Instance.Proxy.Username, config.Instance.Proxy.Password,
		"tcp", config.Instance.Proxy.Address, config.Instance.Proxy.Port, config.Instance.Proxy.Database)
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		logrus.Errorf("[proxy]连接出现错误,url:%v,错误:%v", dsn, e.Error())
		return nil, e
	}
	db.SetConnMaxLifetime(time.Duration(config.Instance.Proxy.LifeTime) * time.Second)
	db.SetMaxOpenConns(config.Instance.Proxy.MaxOpen)
	db.SetMaxIdleConns(config.Instance.Proxy.MaxIdle)
	return &Conn{db}, nil
}

func (c *Conn) Close() {
	if c.DB != nil {
		c.DB.Close()
		c.DB = nil
	}
}

func (c *Conn) Execs(sqlString string, params [][]interface{}) error {
	logrus.Debugf("[proxy]Inserts:%s,params:%v", sqlString, params)
	return c.Tx(func(tx *sql.Tx) error {
		stmt, e := tx.Prepare(sqlString)
		if e != nil {
			return e
		}
		defer stmt.Close()
		for _, value := range params {
			_, e := stmt.Exec(value...)
			if e != nil {
				return e
			}
		}
		return nil
	})
}

func (c *Conn) Exec(sqlString string, param ...interface{}) error {
	logrus.Debugf("[proxy]Insert:%s,param:%v", sqlString, param)
	return c.Tx(func(tx *sql.Tx) error {
		stmt, e := tx.Prepare(sqlString)
		if e != nil {
			return e
		}
		defer stmt.Close()
		_, e = stmt.Exec(param...)
		return e
	})
}

func (c *Conn) Tx(f func(tx *sql.Tx) error) error {
	logrus.Debugf("[proxy]Tx:%v", f)
	tx, e := c.DB.Begin()
	if e != nil {
		return e
	}
	e = f(tx)
	if e != nil {
		defer tx.Rollback()
		return e
	}
	return tx.Commit()
}

func (c *Conn) QueryList(sqlString string, newRow func(types []*sql.ColumnType) []interface{}, param ...interface{}) error {
	return c.Query(sqlString, newRow, func(row []interface{}) {
		//忽略
	}, func(strings []string) {
		//忽略
	}, param...)
}

func (c *Conn) Query(query string, newRow func(types []*sql.ColumnType) []interface{}, rowAfter func(row []interface{}), setColumnNames func([]string), param ...interface{}) error {
	logrus.Debugf("[proxy]Query:%s,param:%v", query, param)
	stmt, e := c.Prepare(query)
	if e != nil {
		return e
	}
	rows, e := stmt.Query(param...)
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
