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
	bytes, e := json.Marshal(event)
	if e != nil {
		return e
	}
	e = s.producer.Send(s.ctx, pulsar.ProducerMessage{
		Payload:    bytes,
		Key:        getKey(event),
		Properties: map[string]string{event.AggType: event.AggId},
	})
	if e != nil {
		logrus.Errorf("[pulsar]发送事件错误,错误:%v", e)
		return e
	} else {
		logrus.Debugf("[pulsar]发送事件成功,聚合[%s-%s],版本[%v]", event.AggType, event.AggId, event.Version)
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
		logrus.Errorf("[pulsar]连接pulsar出现错误,错误:%v", e.Error())
		return nil, e
	}

	producer, e := client.CreateProducer(pulsar.ProducerOptions{
		Topic: config.Instance.Pulsar.TopicName,
	})
	if e != nil {
		logrus.Errorf("[pulsar]创建pulsar.producer出现错误,错误:%v", e.Error())
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

func getKey(event *core.Event) string {
	return fmt.Sprintf("%s:%s", event.AggType, event.AggId)
}
