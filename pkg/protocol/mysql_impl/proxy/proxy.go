package proxy

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/repository"
	"time"
)

var db *sql.DB

func SetConn(conn *repository.Conn) {
	db = conn.DB
}

func Query(query string, newRow func(types []*sql.ColumnType) []interface{}, rowAfter func(row []interface{}), setColumnNames func([]string), param ...interface{}) error {
	logrus.Debugf("[proxy]Query:%s,param:%v", query, param)
	stmt, e := db.Prepare(query)
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

func Proxy(query string) (columnNames []string, columnValues [][]interface{}, err error) {
	logrus.Debugf("[proxy]Mysql:%s", query)
	var temp interface{} = ""
	var rowOrigin []interface{}
	var result []interface{}

	columnValues = make([][]interface{}, 0)
	err = Query(query,
		func(types []*sql.ColumnType) []interface{} {
			if result == nil {
				result = make([]interface{}, len(types))
				rowOrigin = make([]interface{}, 0, len(types))
				for key := range types {
					rowOrigin = append(rowOrigin, temp)
					result[key] = &rowOrigin[key]
				}
			}
			return result
		},
		func(row []interface{}) {
			i := make([]interface{}, len(row))
			for key := range row {
				v1 := rowOrigin[key]
				switch v1.(type) {
				case time.Time:
					i[key] = v1.(time.Time).String()
					continue
				}
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
