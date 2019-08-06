package proxy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"os"
	"testing"
	"time"
)

func TestQueryOne(t *testing.T) {
	config.Instance.Proxy = &config.ProxyConfig{
		Address:  "127.0.0.2",
		Port:     "3306",
		Database: "mysql",
		Username: "root",
		Password: "123456",
		LifeTime: 10,
		MaxOpen:  5,
		MaxIdle:  5,
	}
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=true", config.Instance.Proxy.Username, config.Instance.Proxy.Password,
		"tcp", config.Instance.Proxy.Address, config.Instance.Proxy.Port, config.Instance.Proxy.Database)
	var e error
	db, e = sql.Open("mysql", dsn)
	if e != nil {
		logrus.Errorf("[proxy]连接出现错误,url:%v,错误:%v", dsn, e.Error())
		os.Exit(1)
	}
	defer db.Close()

	var cost_name sql.NullString
	var cost_value sql.NullString
	var last_update time.Time
	var comment sql.NullString
	var default_value sql.NullFloat64
	scan := []interface{}{&cost_name, &cost_value, &last_update, &comment, &default_value}
	err := QueryOne(`select * from server_cost where cost_name like CONCAT('%',?,'%')`, scan, "disk")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(cost_name, cost_value, last_update, comment, default_value)
}

func TestJsonMerge(t *testing.T) {
	var result map[string]interface{}
	var json2 = `{"name":"test","age":18,"name":"test2","age":20,"time":"test","age":null}`
	e := json.Unmarshal([]byte(json2), &result)
	if e != nil {
		panic(e.Error())
	}
	fmt.Println(result)
}
