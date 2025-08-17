package common

import (
	"github.com/rabbitmq/amqp091-go"
)

type Callback func(msgs <-chan amqp091.Delivery)

func ConsumeEvents(q amqp091.Queue, ch *amqp091.Channel, callbackFn Callback) {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	go callbackFn(msgs)
}
