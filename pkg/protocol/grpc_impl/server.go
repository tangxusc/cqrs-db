package rpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"google.golang.org/grpc"
	"net"
	"os"
)

type GrpcServer struct {
	*grpc.Server
	listener net.Listener
}

func NewGrpcServer() (*GrpcServer, error) {
	listener, e := net.Listen(`tcp`, fmt.Sprintf(":%v", config.Instance.Grpc.Port))
	if e != nil {
		return nil, e
	}
	server := grpc.NewServer()

	return &GrpcServer{
		Server:   server,
		listener: listener,
	}, nil
}

func (g *GrpcServer) Start(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			g.Server.GracefulStop()
		}
	}()
	if e := g.Server.Serve(g.listener); e != nil {
		logrus.Errorf("[grpc] start error:%v", e)
		os.Exit(1)
	}
}

func (g *GrpcServer) RegisterService(f func(server *GrpcServer)) {
	f(g)
}
