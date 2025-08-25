package main

import (
	common "common"
	"log"

	// "github.com/davecgh/go-spew/spew"
	"github.com/rabbitmq/amqp091-go"
)

func consomeLances(msgs <-chan amqp091.Delivery) {
	// publicKey, _ := common.ReadAndParseKey("ms-lance/keys/public/1.pem")

	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LANCE: %s", d.Body)

		d.Ack(true)
	}
}

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	qLance, err := common.CreateOrGetQueueAndBind("lance_realizado", ch)
	common.FailOnError(err, "Error connecting to queue")

	common.ConsumeEvents(qLance, ch, consomeLances)

	var forever chan struct{}
	<-forever
}
