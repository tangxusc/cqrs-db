package event

import (
	"fmt"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"runtime"
	"testing"
)

func TestSendEvents(t *testing.T) {
	msgChannel := make(chan pulsar.ConsumerMessage)

	consumerOpts := pulsar.ConsumerOptions{
		Topic:            "cqrs-db",
		SubscriptionName: "my-subscription-1",
		Type:             pulsar.KeyShared,
		MessageChannel:   msgChannel,
	}
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:                     "pulsar://localhost:6650",
		OperationTimeoutSeconds: 5,
		MessageListenerThreads:  runtime.NumCPU(),
	})

	if err != nil {
		panic(err.Error())
		return
	}

	consumer, err := client.Subscribe(consumerOpts)

	if err != nil {
		panic(err.Error())
	}

	defer consumer.Close()

	for cm := range msgChannel {
		msg := cm.Message

		fmt.Printf("收到消息Message ID: %s \n", msg.ID())
		fmt.Printf("收到消息Message value: %s \n", string(msg.Payload()))

		consumer.Ack(msg)
	}
}
