package main

import (
	"log"
	"time"

	common "common"
)

func getLeiloes() string {

	start := time.Now()
	end := time.Now().Add(10000)

	s := "1; Geladeira; " + start.String() + "; " + end.String()

	return s
}

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	qIniciado, err := common.CreateOrGetQueueAndBind("leilao_iniciado", ch)
	common.FailOnError(err, "Error connecting to queue")

	body := getLeiloes()

	common.PublishInQueue(ch, qIniciado, body)

	log.Printf(" [x] Sent %s", body)
}
