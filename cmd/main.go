package main

import (
	"context"
	"cqrs-db/pkg/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	newCommand := cmd.NewCommand(ctx)
	cmd.HandlerNotify(cancel)

	if err := newCommand.Execute(); err != nil {
		logrus.Errorf("发生了错误,错误:%v", err.Error())
	}
}
