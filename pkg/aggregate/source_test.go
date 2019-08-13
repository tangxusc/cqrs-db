package aggregate

import (
	"context"
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"testing"
	"time"
)

/*
无快照
*/
func TestSourcing(t *testing.T) {
	config.Instance.Proxy = &config.ProxyConfig{
		Address:  "127.0.0.2",
		Port:     "3306",
		Database: "test",
		Username: "root",
		Password: "123456",
		LifeTime: 10,
		MaxOpen:  2,
		MaxIdle:  2,
	}
	proxy.InitConn(context.TODO())
	defer proxy.CloseConn()

	//插入事件
	//create table test.event
	//(
	//	id varchar(36) not null,
	//type varchar(20) null,
	//	agg_id varchar(36) null,
	//	agg_type varchar(36) null,
	//	create_time timestamp null,
	//	data varchar(500) null,
	//	constraint event_pk
	//primary key (id)
	//);
	//
	//create table test.snapshot
	//(
	//	id varchar(36) not null,
	//	agg_id varchar(36) null,
	//	agg_type varchar(36) null,
	//	create_time timestamp null,
	//	data varchar(500) null,
	//	constraint snapshot_pk
	//primary key (id)
	//);

	//	insert into test.event
	//	values ('1', 'E1', '1', 'A1', str_to_date('2018-05-02', '%Y-%m-%d %H'), '{"name":"test1"}');
	//insert into test.event
	//values ('2', 'E1', '1', 'A1', str_to_date('2018-05-03', '%Y-%m-%d %H'), '{"age":10}');
	//insert into test.event
	//values ('3', 'E1', '1', 'A1', str_to_date('2018-05-04', '%Y-%m-%d %H'), '{"name":"test2"}');
	//insert into test.event
	//values ('4', 'E1', '1', 'A1', str_to_date('2018-05-05', '%Y-%m-%d %H'), '{"name":"test3","age":null}');

	handler := &db.ConnHandler{
		TxBegin: false,
		TxKey:   "",
	}
	source, _ := GetSource("1", "A1", handler)
	data, err := source.Sourcing(handler)
	fmt.Println("func1", data, err)
	time.Sleep(time.Second * 20)
	fmt.Println("func1 end")

}

/*
基于快照
*/
func TestSourcingWithSnapshot(t *testing.T) {
	config.Instance.Proxy = &config.ProxyConfig{
		Address:  "127.0.0.2",
		Port:     "3306",
		Database: "test",
		Username: "root",
		Password: "123456",
		LifeTime: 10,
		MaxOpen:  2,
		MaxIdle:  2,
	}
	proxy.InitConn(context.TODO())
	defer proxy.CloseConn()

	//插入事件
	//create table test.event
	//(
	//	id varchar(36) not null,
	//type varchar(20) null,
	//	agg_id varchar(36) null,
	//	agg_type varchar(36) null,
	//	create_time timestamp null,
	//	data varchar(500) null,
	//	constraint event_pk
	//primary key (id)
	//);
	//
	//create table test.snapshot
	//(
	//	id varchar(36) not null,
	//	agg_id varchar(36) null,
	//	agg_type varchar(36) null,
	//	create_time timestamp null,
	//	data varchar(500) null,
	//	constraint snapshot_pk
	//primary key (id)
	//);

	//	insert into test.event
	//	values ('1', 'E1', '1', 'A1', str_to_date('2018-05-02', '%Y-%m-%d %H'), '{"name":"test1"}');
	//insert into test.event
	//values ('2', 'E1', '1', 'A1', str_to_date('2018-05-03', '%Y-%m-%d %H'), '{"age":10}');
	//insert into test.event
	//values ('3', 'E1', '1', 'A1', str_to_date('2018-05-04', '%Y-%m-%d %H'), '{"name":"test2"}');
	//insert into test.event
	//values ('4', 'E1', '1', 'A1', str_to_date('2018-05-05', '%Y-%m-%d %H'), '{"name":"test3","age":null}');

	//insert into test.snapshot
	//values ('1', '1', 'A1', str_to_date('2018-05-03', '%Y-%m-%d %H'), '{"name":"test1","age":10}');
	handler := &db.ConnHandler{
		TxBegin: false,
		TxKey:   "",
	}
	source, _ := GetSource("1", "A1", handler)
	data, err := source.Sourcing(handler)
	fmt.Println("func1", data, err)
	time.Sleep(time.Second * 20)
	fmt.Println("func1 end")

}
