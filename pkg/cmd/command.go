package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"github.com/tangxusc/cqrs-db/pkg/db"
	_ "github.com/tangxusc/cqrs-db/pkg/db/handler"
	"github.com/tangxusc/cqrs-db/pkg/event"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func NewCommand(ctx context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:   "start",
		Short: "start db",
		Run: func(cmd *cobra.Command, args []string) {
			rand.Seed(time.Now().Unix())
			config.InitLog()
			//启动数据库
			go db.Start(ctx)
			//连接代理数据库
			go proxy.InitConn(ctx)
			defer proxy.CloseConn()
			//启动事件恢复机制
			go event.RecoveryEvent(ctx)
			defer event.Stop()
			//启动事件发送
			go event.Start(ctx)
			defer event.Close()

			<-ctx.Done()
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
