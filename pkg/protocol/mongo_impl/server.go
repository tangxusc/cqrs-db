package mongo_impl

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	protocol "github.com/tangxusc/mongo-protocol"
	"os"
)

func (s *MongoServer) Start(ctx context.Context) {
	if e := s.Server.Start(ctx); e != nil {
		logrus.Errorf(`[mongo]server error: %v`, e)
		os.Exit(1)
	}
}

type MongoServer struct {
	*protocol.Server
	Port                string
	QueryHandlers       []protocol.QueryHandler
	DefaultQueryHandler protocol.QueryHandler
}

func (s *MongoServer) Process(header *protocol.MsgHeader, r *protocol.Reader, conn *protocol.ConnContext) error {
	query := &protocol.Query{}
	query.Header = *header
	e := query.UnMarshal(r)
	if e != nil {
		return e
	}
	return s.DoQueryHandler(header, conn, query)
}

func (s *MongoServer) DoQueryHandler(header *protocol.MsgHeader, conn *protocol.ConnContext, query *protocol.Query) error {
	reply := protocol.NewReply(header.RequestID)
	defer func() {
		if e := recover(); e != nil {
			reply.ResponseFlags = protocol.QueryFailure
			reply.NumberReturned = 1
			reply.Documents = map[string]interface{}{"$err": fmt.Sprintf("%v", e)}
		}
		if e := reply.Write(conn); e != nil {
			panic(e)
		}
	}()
	for _, item := range s.QueryHandlers {
		if item.Support(query) {
			return item.Process(query, reply)
		}
	}
	return s.DefaultQueryHandler.Process(query, reply)
}

func NewMongoServer(port string) *MongoServer {
	newServer := protocol.NewServer(port)
	server := &MongoServer{
		Server: newServer,
		Port:   port,
	}
	server.DefaultQueryHandler = &defaultQueryHandler{}
	server.AddHandler(protocol.OP_QUERY, server)
	return server
}

func (s *MongoServer) AddQueryHandler(handlers ...protocol.QueryHandler) {
	s.QueryHandlers = append(s.QueryHandlers, handlers...)
}

type defaultQueryHandler struct {
}

func (d *defaultQueryHandler) Support(query *protocol.Query) bool {
	return true
}

func (d *defaultQueryHandler) Process(query *protocol.Query, reply *protocol.Reply) error {
	reply.NumberReturned = 1
	reply.Documents = map[string]interface{}{"ok": 1}
	return nil
}
