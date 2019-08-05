package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	newCommand := cmd.NewCommand(ctx)
	cmd.HandlerNotify(cancel)

	if err := newCommand.Execute(); err != nil {
		logrus.Errorf("发生了错误,错误:%v", err.Error())
	}
}
