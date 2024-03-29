package mysql_impl

import (
	"crypto/tls"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/siddontang/go-mysql/test_util/test_keys"
	"net"
	"testing"
	"time"
)

type RemoteThrottleProvider struct {
	*server.InMemoryProvider
	delay int // in milliseconds
}

func (m *RemoteThrottleProvider) GetCredential(username string) (password string, found bool, err error) {
	time.Sleep(time.Millisecond * time.Duration(m.delay))
	return m.InMemoryProvider.GetCredential(username)
}

func TestStart(t *testing.T) {
	l, _ := net.Listen("tcp", "127.0.0.1:3308")
	// user either the in-memory credential provider or the remote credential provider (you can implement your own)
	//inMemProvider := server.NewInMemoryProvider()
	//inMemProvider.AddUser("root", "123")
	remoteProvider := &RemoteThrottleProvider{server.NewInMemoryProvider(), 10 + 50}
	remoteProvider.AddUser("root", "123")
	var tlsConf = server.NewServerTLSConfig(test_keys.CaPem, test_keys.CertPem, test_keys.KeyPem, tls.VerifyClientCertIfGiven)
	for {
		c, _ := l.Accept()
		go func() {
			// Create a connection with user root and an empty password.
			// You can use your own handler to handle command here.
			svr := server.NewServer("8.0.12", mysql.DEFAULT_COLLATION_ID, mysql.AUTH_SHA256_PASSWORD, test_keys.PubPem, tlsConf)
			conn, err := server.NewCustomizedConn(c, svr, remoteProvider, server.EmptyHandler{})

			if err != nil {
				fmt.Println("Connection error:", err)
				return
			}

			for {
				conn.HandleCommand()
			}
		}()
	}
}

func TestAbcd(t *testing.T) {

}
