package main

import (
	"context"
	"github.com/tangxusc/cqrs-db/pkg/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	newCommand := cmd.NewCommand(ctx)
	cmd.HandlerNotify(cancel)

	_ = newCommand.Execute()
	cancel()
}
