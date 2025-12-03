package common

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToBroker() (*amqp.Connection, *amqp.Channel) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		EXCHANGE_NAME, // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	FailOnError(err, "Failed to declare an exchange")

	return conn, ch
}
