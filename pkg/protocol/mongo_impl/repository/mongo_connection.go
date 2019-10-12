package repository

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type MongoConn struct {
	*mongo.Client
	dbName string
}

func (c *MongoConn) Conn(ctx context.Context) (e error) {
	//https://github.com/hwholiday/learning_tools/blob/master/mongodb/mongo-go-driver/main.go
	opt := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.Instance.Mongo.Address, config.Instance.Mongo.Port))
	opt.SetLocalThreshold(time.Duration(config.Instance.Mongo.LocalThreshold) * time.Second)
	opt.SetMaxConnIdleTime(time.Duration(config.Instance.Mongo.MaxConnIdleTime) * time.Second)
	opt.SetMaxPoolSize(uint64(config.Instance.Mongo.MaxPoolSize))
	//表示只使用辅助节点
	//opt.SetReadPreference(want)
	//指定查询应返回实例的最新数据确认为，已写入副本集中的大多数成员
	opt.SetReadConcern(readconcern.Majority())
	//请求确认写操作传播到大多数mongod实例
	opt.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	opt.SetAuth(options.Credential{
		Username: config.Instance.Mongo.Username,
		Password: config.Instance.Mongo.Password,
	})
	c.Client, e = mongo.Connect(ctx, opt)
	if e != nil {
		logrus.Errorf("[db]connection mongodb error:%v", e)
		return
	}
	if e = c.Client.Ping(ctx, readpref.Primary()); e != nil {
		logrus.Errorf("[db]ping mongodb error:%v", e)
		return
	}
	c.dbName = config.Instance.Mongo.DbName
	return
}

func (c *MongoConn) Close(ctx context.Context) {
	e := c.Client.Disconnect(ctx)
	if e != nil {
		logrus.Errorf("[db]disconnection mongodb error:%v", e)
	}
}
