package main

import (
	common "common"
	"fmt"
	"log"

	// "github.com/davecgh/go-spew/spew"
	// "github.com/davecgh/go-spew/spew"
	"github.com/rabbitmq/amqp091-go"
)

var connIn *amqp091.Connection
var chIn *amqp091.Channel

var connOut *amqp091.Connection
var chOut *amqp091.Channel

func handleLanceValidado(lanceValidad []byte) {
	log.Printf("Novo lance validado: ")
	lance := common.ByteArrayToLance(lanceValidad)
	routing_key := lance.LeilaoID

	nome_fila := fmt.Sprintf("leilao_%s", routing_key)

	notificacao := common.Notificacao{
		Lance:  lance,
		Status: common.NovoLance,
	}

	byteNotificacao := common.NotificacaoToByteArray(notificacao)

	q, err := common.CreateOrGetQueueAndBind(nome_fila, chOut)
	common.FailOnError(err, "Error connecting to queue")

	common.PublishInQueue(chOut, q, byteNotificacao, nome_fila)

}

func consomeLances(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LANCE VALIDADO: %s", d.Body)

		go handleLanceValidado(d.Body)

		d.Ack(false)
	}
}

func handleLanceGanhador(lanceGanhador []byte) {
	log.Printf("Novo lance ganhador: ")
	lance := common.ByteArrayToLance(lanceGanhador)
	routing_key := lance.LeilaoID

	nome_fila := fmt.Sprintf("leilao_%s", routing_key)

	notificacao := common.Notificacao{
		Lance:  lance,
		Status: common.GanhadorLance,
	}

	byteNotificacao := common.NotificacaoToByteArray(notificacao)

	q, err := common.CreateOrGetQueueAndBind(nome_fila, chOut)
	common.FailOnError(err, "Error connecting to queue")

	common.PublishInQueue(chOut, q, byteNotificacao, nome_fila)

}

func consomeLancesGanhador(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LANCE GANHADOR: %s", d.Body)

		go handleLanceGanhador(d.Body)

		d.Ack(false)
	}
}

func main() {
	connIn, chIn = common.ConnectToBroker()
	defer connIn.Close()
	defer chIn.Close()

	connOut, chOut = common.ConnectToBroker()
	defer connOut.Close()
	defer chOut.Close()

	qLanceVal, err := common.CreateOrGetQueueAndBind("lance_validado", chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLanceVal, chIn, consomeLances)

	qLanceWin, err := common.CreateOrGetQueueAndBind("leilao_vencedor", chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLanceWin, chIn, consomeLancesGanhador)

	var forever chan struct{}
	<-forever
}
