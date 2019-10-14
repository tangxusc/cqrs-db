package mongo_impl

import (
	"context"
	"github.com/sirupsen/logrus"
	protocol "github.com/tangxusc/mongo-protocol"
)

func (s *MongoServer) Start(ctx context.Context) {
	if e := s.Server.Start(ctx); e != nil {
		logrus.Errorf(`[mongo]server error: %v`, e)
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
	e := query.UnMarshal(r)
	if e != nil {
		return e
	}
	reply := protocol.NewReply(header.RequestID)
	defer func() {
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
	server := &MongoServer{
		Server: protocol.NewServer(port),
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
