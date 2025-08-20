package common

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishInQueue(ch *amqp.Channel, q amqp.Queue, message []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(ctx,
		"LEILAO", // exchange
		q.Name,   // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	FailOnError(err, "Failed to publish a message")
}
