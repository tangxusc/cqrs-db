package mysql_impl

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestTestHandler_HandleQuery(t *testing.T) {
	var test1 string
	var test2 string

	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=true", "root", "",
		"tcp", "127.0.0.1", "3307", "test")
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		panic(e.Error())
	}
	defer db.Close()
	tx, e := db.Begin()
	if e != nil {
		panic(e.Error())
	}
	defer tx.Commit()
	stmt, e := tx.Prepare("select * from some_table_name")
	if e != nil {
		panic(e.Error())
	}
	defer stmt.Close()
	rows, e := stmt.Query()
	if e != nil {
		panic(e.Error())
	}
	defer rows.Close()
	strings, e := rows.Columns()
	fmt.Println("获取的行:", strings, e)
	for rows.Next() {
		e := rows.Scan(&test1, &test2)
		if e != nil {
			panic(e.Error())
		}
		fmt.Println("test1:", test1, ",test2:", test2)
	}

}
