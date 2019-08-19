package event

import (
	"context"
	"fmt"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/cqrs-db/pkg/config"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
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
		e := p.sendEvent(value)
		if e != nil {
			logrus.Errorf("[event]发送事件错误,错误:%v", e)
			return e
		}
	}
	return nil
}

func (p *PulsarSender) sendEvent(event Event) error {
	bytes, e := event.ToJson()
	if e != nil {
		return e
	}
	p.producer.SendAsync(p.ctx, pulsar.ProducerMessage{
		Payload:    bytes,
		Key:        getKey(event),
		Properties: map[string]string{event.AggType(): event.AggId()},
	}, func(message pulsar.ProducerMessage, e error) {
		if e != nil {
			logrus.Errorf("[event]发送事件错误,错误:%v", e)
			return
		}
		for k, v := range message.Properties {
			e := proxy.Exec(`update event set mq_status=? where agg_id=? and agg_type=? and mq_status=?`, NotSend, v, k, Sent)
			if e != nil {
				logrus.Errorf("[event]事件已发送,更新数据库发生错误:%s,%s,错误:%v", k, v, e)
			}
		}
	})
	return e
}

func getKey(event Event) string {
	return fmt.Sprintf("%s:%s", event.AggType(), event.AggId())
}
