package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/memory"
	"github.com/tangxusc/cqrs-db/pkg/mq"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mongo_impl"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mongo_impl/handler"
	mongoRepository "github.com/tangxusc/cqrs-db/pkg/protocol/mongo_impl/repository"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl"
	_ "github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/handler"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/proxy"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/repository"
	protocol "github.com/tangxusc/mongo-protocol"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	newCommand := NewCommand(ctx)
	HandlerNotify(cancel)

	_ = newCommand.Execute()
	cancel()
}

func NewCommand(ctx context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:   "start",
		Short: "start db",
		RunE: func(cmd *cobra.Command, args []string) error {
			rand.Seed(time.Now().Unix())
			config.InitLog()

			if e := StartMysqlProtocol(ctx); e != nil {
				return e
			}

			if e := StartMongoProtocol(ctx); e != nil {
				return e
			}

			core.SetSnapshotSaveStrategyFactory(repository.NewCountStrategyFactory(config.Instance.ServerDb.MaxEventToSnapshot))
			impl := memory.NewAggregateManagerImpl(ctx)
			core.SetAggregateManager(impl)

			if config.Instance.Pulsar.Enable {
				sender, e := mq.NewSender(ctx)
				if e != nil {
					return e
				}
				defer sender.Close()
				core.SetEventSender(sender)
			}

			//启动事件恢复机制
			core.NewRestorer().Start(ctx)

			<-ctx.Done()
			return nil
		},
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	config.BindParameter(command)

	return command
}

func StartMongoProtocol(ctx context.Context) error {
	if config.Instance.Mongo.Enable {
		mongo := mongo_impl.NewMongoServer(config.Instance.ServerDb.Port)
		mongo.AddQueryHandler(handler.GetBaseQueryHandler()...)
		mongo.AddQueryHandler(handler.NewFindHandler())
		insertHandler := handler.NewInsertHandler()
		mongo.AddHandler(protocol.OP_INSERT, insertHandler)
		mongo.AddQueryHandler(handler.NewInsertFindHandler(insertHandler))
		conn := mongoRepository.NewMongoConn()
		e := conn.Conn(ctx)
		if e != nil {
			return e
		}
		core.SetEventStore(mongoRepository.NewEventStoreImpl(conn, ctx, config.Instance.Mongo.EventCollectionName))
		core.SetSnapshotStore(mongoRepository.NewSnapshotStoreImpl(conn, ctx, config.Instance.Mongo.SnapshotCollectionName))
		go mongo.Start(ctx)
	}
	return nil
}

func StartMysqlProtocol(ctx context.Context) error {
	if config.Instance.Mysql.Enable {
		//连接代理数据库
		conn, e := repository.InitConn(ctx)
		if e != nil {
			return e
		}
		//启动mysql协议
		go mysql_impl.Start(ctx)
		proxy.SetConn(conn)
		impl := repository.NewEventStoreImpl(conn)
		core.SetEventStore(impl)
		core.SetSnapshotStore(repository.NewSnapshotStoreImpl(conn))
		return e
	}
	return nil
}

func HandlerNotify(cancel context.CancelFunc) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		cancel()
	}()
}
