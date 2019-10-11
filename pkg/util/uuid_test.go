package util

import (
	"context"
	"fmt"
	protocol "github.com/tangxusc/mongo-protocol"
	"testing"
)

func TestGenerateUuid(t *testing.T) {
	uuid := GenerateUuid()
	fmt.Println(uuid)
}

func TestMongoServer(t *testing.T) {
	server := protocol.NewServer(`27017`)
	e := server.Start(context.TODO())
	if e != nil {
		panic(e)
	}
}
