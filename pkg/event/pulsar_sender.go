package event

import (
	"context"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"os"
	"runtime"
)

type PulsarSender struct {
	client   pulsar.Client
	producer pulsar.Producer
	ctx      context.Context
}

func (p *PulsarSender) Close() error {
	_ = p.producer.Close()
	return p.client.Close()
}

func Start(ctx context.Context) {
	if len(config.Instance.Pulsar.Url) <= 0 {
		return
	}
	var err error
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:                     config.Instance.Pulsar.Url,
		OperationTimeoutSeconds: 5,
		MessageListenerThreads:  runtime.NumCPU(),
	})
	if err != nil {
		logrus.Errorf("[event]连接pulsar出现错误,错误:%v", err.Error())
		os.Exit(1)
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: config.Instance.Pulsar.TopicName,
	})
	if err != nil {
		logrus.Errorf("[event]创建pulsar.producer出现错误,错误:%v", err.Error())
		os.Exit(1)
	}
	SenderImpl = &PulsarSender{
		client:   client,
		producer: producer,
		ctx:      ctx,
	}
}

func Close() {
	sender, ok := SenderImpl.(*PulsarSender)
	if !ok {
		return
	}
	_ = sender.Close()
}

func (p *PulsarSender) SendEvents(events []Event) error {
	for _, value := range events {
		p.sendEvent(value)
	}
	return nil
}

func (p *PulsarSender) sendEvent(event Event) {
	p.producer.SendAsync(p.ctx, pulsar.ProducerMessage{
		//todo:实现序列化
		Payload: []byte{},
		Value:   event,
	}, func(message pulsar.ProducerMessage, e error) {
		//todo:实现更新数据库状态
	})
}
