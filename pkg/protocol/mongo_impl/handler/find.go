package handler

import (
	protocol "github.com/tangxusc/mongo-protocol"
)

type Handler struct {
}

func (q *Handler) Process(header *protocol.MsgHeader, r *protocol.Reader, conn *protocol.ConnContext) error {
	return nil
}
