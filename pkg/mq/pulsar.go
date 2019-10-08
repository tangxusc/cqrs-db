package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"runtime"
)

type EventSenderImpl struct {
	client   pulsar.Client
	producer pulsar.Producer
	ctx      context.Context
}

func (s *EventSenderImpl) Send(event *core.Event) error {
	e := s.sendEvent(event)
	if e != nil {
		logrus.Errorf("[event]发送事件错误,错误:%v", e)
		return e
	}
	return nil
}

func NewSender(ctx context.Context) (sender *EventSenderImpl, e error) {
	if len(config.Instance.Pulsar.Url) <= 0 {
		return nil, fmt.Errorf("pulsar.url not set")
	}
	client, e := pulsar.NewClient(pulsar.ClientOptions{
		URL:                     config.Instance.Pulsar.Url,
		OperationTimeoutSeconds: 5,
		MessageListenerThreads:  runtime.NumCPU(),
	})
	if e != nil {
		logrus.Errorf("[event]连接pulsar出现错误,错误:%v", e.Error())
		return nil, e
	}

	producer, e := client.CreateProducer(pulsar.ProducerOptions{
		Topic: config.Instance.Pulsar.TopicName,
	})
	if e != nil {
		logrus.Errorf("[event]创建pulsar.producer出现错误,错误:%v", e.Error())
		return nil, e
	}
	sender = &EventSenderImpl{
		client:   client,
		producer: producer,
		ctx:      ctx,
	}
	return
}

func (s *EventSenderImpl) Close() {
	s.producer.Close()
	s.client.Close()
}

func (s *EventSenderImpl) sendEvent(event *core.Event) error {
	bytes, e := json.Marshal(event)
	if e != nil {
		return e
	}
	return s.producer.Send(s.ctx, pulsar.ProducerMessage{
		Payload:    bytes,
		Key:        getKey(event),
		Properties: map[string]string{event.AggType: event.AggId},
	})
}

func getKey(event *core.Event) string {
	return fmt.Sprintf("%s:%s", event.AggType, event.AggId)
}
