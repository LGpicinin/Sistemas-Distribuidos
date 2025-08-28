package main

import (
	common "common"
)

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	// qLeilao1, err := common.CreateOrGetQueueAndBind("leilao_1", ch)
	// common.FailOnError(err, "Error connecting to queue")

	// var b []byte

	// common.PublishInQueue(ch, qLeilao1, b)
}
