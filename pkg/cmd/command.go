package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"github.com/tangxusc/cqrs-db/pkg/db"
	_ "github.com/tangxusc/cqrs-db/pkg/db/handler"
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

			go db.Start(ctx)
			go proxy.InitConn(ctx)
			defer proxy.CloseConn()

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
