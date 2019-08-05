package db

import (
	"context"
	"cqrs-db/pkg/config"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/sirupsen/logrus"
	"net"
	"os"
)

const serverVersion = "8.0.3"

func Start(ctx context.Context) {
	l, e := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", config.Instance.Db.Port))
	if e != nil {
		logrus.Errorf("[db]监听tcp出现错误,错误:%v", e.Error())
		os.Exit(1)
	}

	newServer := server.NewServer(serverVersion, mysql.DEFAULT_COLLATION_ID, mysql.AUTH_NATIVE_PASSWORD, nil, nil)
	provider := server.NewInMemoryProvider()
	//TODO:在java模式下,提示密码错误,sha1加密问题?
	provider.AddUser(config.Instance.Db.Username, config.Instance.Db.Password)
	for {
		select {
		case <-ctx.Done():
			e = l.Close()
			return
		default:
			lConn, e := l.Accept()
			if e != nil {
				logrus.Errorf("[db]tcp Accept出现错误,错误:%v", e.Error())
				return
			}
			conn, e := server.NewCustomizedConn(lConn, newServer, provider, NewHandler())
			if e != nil {
				logrus.Errorf("[db]NewCustomizedConn出现错误,错误:%v", e.Error())
				break
			}
			go startConnCommandHandler(ctx, conn)
		}
	}
}

func startConnCommandHandler(ctx context.Context, conn *server.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			e := conn.HandleCommand()
			if e != nil {
				logrus.Warningf("[db]HandleCommand出现错误,错误:%v", e.Error())
				return
			}
		}
	}
}
