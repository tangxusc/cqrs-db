# cqrs-db 

## 介绍

cqrs-db通过代理mysql的协议,实现在proxy层中完成cqrs中的事件溯源,事件存储,事件发送到mq等功能,现处于早期开发阶段.

使用者通过mysql协议连接proxy,proxy解析sql,并将查询特定的表的规则解析为聚合的事件溯源,并通过mysql的json格式返回至客户端

## 运行

### 1.启动mysql&pulsar

```shell
docker run --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql
docker run --rm -it -p 6650:6650 -p 8080:8080 apachepulsar/pulsar:2.3.2 bin/pulsar standalone
docker run -itd --rm -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123456 mongo
```

### 2.初始化mysql

```sql
# 插入事件
create schema test;
create table test.event
(
    id          varchar(36)  not null,
    type        varchar(20)  null,
    agg_id      varchar(36)  null,
    agg_type    varchar(36)  null,
    create_time timestamp    null,
    data        varchar(500) null,
    mq_status   varchar(50)  null,
    revision    int,
    constraint event_pk
        primary key (id)
);

create table test.snapshot
(
    id          varchar(36)  not null,
    agg_id      varchar(36)  null,
    agg_type    varchar(36)  null,
    create_time timestamp    null,
    data        varchar(500) null,
    revision    int,
    constraint snapshot_pk
        primary key (id)
);

insert into test.event
values ('1', 'E1', '1', 'A1', str_to_date('2018-05-02', '%Y-%m-%d %H'), '{"name":"test1"}', 'NotSend',1);
insert into test.event
values ('2', 'E1', '1', 'A1', str_to_date('2018-05-03', '%Y-%m-%d %H'), '{"age":10}', 'NotSend',2);
insert into test.event
values ('3', 'E1', '1', 'A1', str_to_date('2018-05-04', '%Y-%m-%d %H'), '{"name":"test2"}', 'NotSend',3);
insert into test.event
values ('4', 'E1', '1', 'A1', str_to_date('2018-05-05', '%Y-%m-%d %H'), '{"name":"test3","age":null}', 'NotSend',4);

insert into test.snapshot
values ('1', '1', 'A1', str_to_date('2018-05-03', '%Y-%m-%d %H'), '{"name":"test1","age":10}',2);
```

在mysql中初始化部分测试数据,方便我们使用

### 3.运行proxy

```shell
cd cqrs-db/cmd/
go run main.go --debug true --proxy-Database=test --proxy-Password=123456 --proxy-Username =root --proxy-address=127.0.0.1 --proxy-port=3306
```

其他参数可以通过`go run main.go --help`查看

### 4.使用mysql客户端连接proxy

```shell
#默认proxy启动在3307端口,默认用户名root,密码:123456
mysql -uroot -P3307 -h 127.0.0.1 -p123456
```

### 5.proxy支持的查询

```sql
#查询聚合类型为a1的聚合,聚合的类型通过表名称指定 xxx_aggregate
begin;
select id as c1, data as c3, agg_type as c2 from a1_aggregate a1 where id = '1';
commit;
#查询聚合的锁
select * from locks_agg;
#查询聚合
select * from agg_info;
#保存事件
begin;
insert into event_aggregate(type, agg_id, agg_type, create_time, data) values ('E1', '1', 'A1', '2018-08-02 12:00:00', '{"name":"test1"}'),('E1', '1', 'A1', '2018-08-03 14:00:00', '{"name":"test5"}');
commit;
```

```sql
db.a1_aggregate.find({'id':'4'})
db.a1_event.insert({'aggId':'1','eventType':'test1','data':{'name':'testname1'},'version':1,'createTime':'2006-01-02 15:04:05'})
```

### 6.mysql本身的表查询

```sql
select * from event;
select * from snapshot;
```

## 参照

```
github.com/go-sql-driver/mysql
github.com/gofrs/uuid
github.com/siddontang/go-mysql
github.com/sirupsen/logrus
github.com/spf13/cobra
github.com/spf13/viper
github.com/xwb1989/sqlparser
github.com/apache/pulsar/pulsar-client-go
```

