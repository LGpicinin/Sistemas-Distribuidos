package main

import (
	common "common"
	"crypto"
	"crypto/rsa"
	"fmt"
	"log"

	// "github.com/davecgh/go-spew/spew"
	// "github.com/davecgh/go-spew/spew"
	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var ch *amqp091.Channel

var activeLeiloes map[string]common.ActiveLeilao

func verifySignature(signedLance common.SignedLance) (bool, error) {
	hashedLance := common.HashLance(signedLance.Lance)

	publicKey, _ := common.ReadAndParseKey(fmt.Sprintf("ms-lance/keys/public/%s.pem", signedLance.Lance.UserID))

	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedLance, signedLance.Signature)
	return err == nil, err
}

func handleLanceCandidate(lanceCanditate []byte) {
	signedLance := common.ByteArrayToSignedLance(lanceCanditate)
	isValidSignature, err := verifySignature(signedLance)
	if !isValidSignature {
		log.Panicf("[MS-LANCE] Erro ao verificar chave: %v", err)
		return
	}

	activeLeilao, ok := activeLeiloes[signedLance.Lance.LeilaoID]
	if !ok {
		return
	}

	if signedLance.Lance.Value <= activeLeilao.LastValidLance.Value {
		return
	}

	activeLeilao.LastValidLance = signedLance.Lance

	q, err := common.CreateOrGetQueueAndBind("lance_validado", ch)
	common.FailOnError(err, "Error connecting to queue")

	common.PublishInQueue(ch, q, common.LanceToByteArray(signedLance.Lance))
}

func consomeLances(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LANCE: %s", d.Body)

		go handleLanceCandidate(d.Body)

		d.Ack(true)
	}
}

func handleLeilaoIniciado(leilaoByteArray []byte) {
	leilao := common.ByteArrayToLeilao(leilaoByteArray)

	activeLeiloes[leilao.ID] = common.ActiveLeilao{
		Leilao: leilao,
	}
}

func consumeLeiloesIniciados(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LEILAO INICIADO: %s", d.Body)

		go handleLeilaoIniciado(d.Body)

		d.Ack(true)
	}
}

func handleLeilaoFinalizado(leilaoByteArray []byte) {
	leilao := common.ByteArrayToLeilao(leilaoByteArray)
	activeLeilao, ok := activeLeiloes[leilao.ID]
	if ok {
		lastLance := activeLeilao.LastValidLance
		if lastLance != (common.Lance{}) {
			q, err := common.CreateOrGetQueueAndBind("leilao_vencedor", ch)
			common.FailOnError(err, "Error connecting to queue")

			common.PublishInQueue(ch, q, common.LanceToByteArray(lastLance))
		}

		delete(activeLeiloes, leilao.ID)
	}
}

func consumeLeiloesFinalizados(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LEILAO FINALIZADO: %s", d.Body)

		go handleLeilaoFinalizado(d.Body)

		d.Ack(true)
	}
}

func main() {
	conn, ch = common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	activeLeiloes = make(map[string]common.ActiveLeilao)

	qLance, err := common.CreateOrGetQueueAndBind("lance_realizado", ch)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLance, ch, consomeLances)

	qLeiloesIniciados, err := common.CreateOrGetQueueAndBind("leilao_iniciado", ch)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesIniciados, ch, consumeLeiloesIniciados)

	qLeiloesFinalizados, err := common.CreateOrGetQueueAndBind("leilao_finalizado", ch)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesFinalizados, ch, consumeLeiloesFinalizados)

	var forever chan struct{}
	<-forever
}
