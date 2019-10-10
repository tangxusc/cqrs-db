package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"github.com/tangxusc/cqrs-db/pkg/memory"
	"github.com/tangxusc/cqrs-db/pkg/mq"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl"
	_ "github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/handler"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/proxy"
	"github.com/tangxusc/cqrs-db/pkg/protocol/mysql_impl/repository"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func NewCommand(ctx context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:   "start",
		Short: "start db",
		RunE: func(cmd *cobra.Command, args []string) error {
			rand.Seed(time.Now().Unix())
			config.InitLog()
			//连接代理数据库
			conn, e := repository.InitConn(ctx)
			if e != nil {
				return e
			}
			defer conn.Close()

			//启动mysql协议
			go mysql_impl.Start(ctx)
			proxy.SetConn(conn)

			core.SetEventStore(repository.NewEventStoreImpl(conn))
			core.SetSnapshotStore(repository.NewSnapshotStoreImpl(conn))
			core.SetSnapshotSaveStrategy(repository.NewCountStrategy(100))
			impl := memory.NewAggregateManagerImpl(ctx)
			core.SetAggregateManager(impl)

			sender, e := mq.NewSender(ctx)
			if e != nil {
				return e
			}
			defer sender.Close()
			core.SetEventSender(sender)

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

func HandlerNotify(cancel context.CancelFunc) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		cancel()
	}()
}
