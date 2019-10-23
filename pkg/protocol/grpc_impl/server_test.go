package rpc

import (
	"context"
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNewGrpcServer(t *testing.T) {
	todo, cancelFunc := context.WithTimeout(context.TODO(), time.Second*100)

	config.Instance.Grpc.Port = `6666`
	conn, e := grpc.Dial(fmt.Sprintf(`:%v`, config.Instance.Grpc.Port), grpc.WithInsecure())
	if e != nil {
		panic(e)
	}
	eventsClient := NewEventsClient(conn)

	go func() {
		for i := 1; i <= 10; i++ {
			var version = int32(i)
			publishRequest := &PublishRequest{
				AggId:     "1",
				AggType:   "A1",
				Version:   version,
				EventType: "E1",
				Data:      fmt.Sprintf(`{"name":"name-%v"}`, i),
			}
			publishResponse, e := eventsClient.Publish(todo, publishRequest)
			if e != nil {
				panic(e)
			}
			fmt.Println(publishResponse.Version)
			time.Sleep(time.Second * 2)
		}
		cancelFunc()
	}()

	go func() {
		for {
			select {
			case <-todo.Done():
				return
			default:
				fmt.Println("==============sourcing=======================")
				sourcingClient := NewSourcingClient(conn)
				request := &SourcingRequest{
					AggId:   "1",
					AggType: "A1",
				}
				response, e := sourcingClient.Sourcing(todo, request)
				if e != nil {
					panic(e)
				}
				fmt.Println(response.Data, `===================================`)
				time.Sleep(time.Second)
			}
		}
	}()

	select {
	case <-todo.Done():
		_ = conn.Close()
		fmt.Println("end")
	}
}
