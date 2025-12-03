package common

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateQueue(queueName string, ch *amqp.Channel) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName,    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return q, err
}

func CreateOrGetQueueAndBind(queueName string, routingKey string, ch *amqp.Channel) (amqp.Queue, error) {
	q, err := CreateQueue(queueName, ch)
	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,        // queue name
		routingKey,     // routing key
		EXCHANGE_NAME, // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind a queue")

	return q, err
}
