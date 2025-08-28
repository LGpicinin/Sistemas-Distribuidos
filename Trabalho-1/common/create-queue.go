package common

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateQueue(ch *amqp.Channel) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return q, err
}

func CreateOrGetQueueAndBind(queueName string, ch *amqp.Channel) (amqp.Queue, error) {
	q, err := CreateQueue(ch)
	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,        // queue name
		queueName,     // routing key
		EXCHANGE_NAME, // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind a queue")

	return q, err
}
