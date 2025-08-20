package main

import (
	common "common"
)

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	qLance, err := common.CreateOrGetQueueAndBind("lance_realizado", ch)
	common.FailOnError(err, "Error connecting to queue")

	var b []byte

	common.PublishInQueue(ch, qLance, b)
}
