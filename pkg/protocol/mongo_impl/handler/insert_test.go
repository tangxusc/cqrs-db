package handler

import (
	"fmt"
	protocol "github.com/tangxusc/mongo-protocol"
	"testing"
)

func TestCheckColumns(t *testing.T) {
	handler := NewInsertHandler()
	insert := &protocol.Insert{
		Flags:              0,
		FullCollectionName: "aaa_event",
		Documents:          nil,
	}
	e := handler.validateDocuments(insert)
	fmt.Println(e)
}
