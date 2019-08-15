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
    constraint snapshot_pk
        primary key (id)
);

insert into test.event
values ('1', 'E1', '1', 'A1', str_to_date('2018-05-02', '%Y-%m-%d %H'), '{"name":"test1"}', 'NotSend');
insert into test.event
values ('2', 'E1', '1', 'A1', str_to_date('2018-05-03', '%Y-%m-%d %H'), '{"age":10}', 'NotSend');
insert into test.event
values ('3', 'E1', '1', 'A1', str_to_date('2018-05-04', '%Y-%m-%d %H'), '{"name":"test2"}', 'NotSend');
insert into test.event
values ('4', 'E1', '1', 'A1', str_to_date('2018-05-05', '%Y-%m-%d %H'), '{"name":"test3","age":null}', 'NotSend');

insert into test.snapshot
values ('1', '1', 'A1', str_to_date('2018-05-03', '%Y-%m-%d %H'), '{"name":"test1","age":10}');

# 查询聚合 mysql -uroot -P3307 -h 127.0.0.1
select id as c1, data as c3, agg_type as c2
from a1_aggregate a1
where id = '1';
# +------+------+-----------------------------+
# | c1   | c2   | c3                          |
# +------+------+-----------------------------+
# | 1    | a1   | {"age":null,"name":"test3"} |
# +------+------+-----------------------------+
# 1 row in set (0.00 sec)