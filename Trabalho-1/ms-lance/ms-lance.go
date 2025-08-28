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

var connIn *amqp091.Connection
var chIn *amqp091.Channel

var connOut *amqp091.Connection
var chOut *amqp091.Channel

var activeLeiloes map[string]common.ActiveLeilao

func verifySignature(signedLance common.SignedLance) (bool, error) {
	hashedLance := common.HashLance(signedLance.Lance)

	publicKey, _ := common.ReadAndParseKey(fmt.Sprintf("ms-lance/keys/public/%s.pem", signedLance.Lance.UserID))

	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedLance, signedLance.Signature)
	return err == nil, err
}

func handleLanceCandidate(lanceCanditate []byte) {
	log.Printf("Novo lance validado:gsadgsajdbs ")
	signedLance := common.ByteArrayToSignedLance(lanceCanditate)
	isValidSignature, err := verifySignature(signedLance)
	if !isValidSignature {
		log.Panicf("[MS-LANCE] Erro ao verificar chave: %v", err)
		return
	}

	activeLeilao, ok := activeLeiloes[signedLance.Lance.LeilaoID]
	if !ok {
		log.Panicf("[MS-LANCE] Erro ao acessar leilão ativo %v", signedLance.Lance.LeilaoID)
		return
	}

	if signedLance.Lance.Value <= activeLeilao.LastValidLance.Value {
		log.Panicf("[MS-LANCE] Lance não válido %v", signedLance.Lance)
		return
	}

	activeLeilao.LastValidLance = signedLance.Lance

	q, err := common.CreateOrGetQueueAndBind("lance_validado", chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.PublishInQueue(chOut, q, common.LanceToByteArray(signedLance.Lance))

	log.Printf("Novo lance validado: %v", signedLance.Lance)
}

func consomeLances(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LANCE: %s", d.Body)

		go handleLanceCandidate(d.Body)

		d.Ack(false)
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

		d.Ack(false)
	}
}

func handleLeilaoFinalizado(leilaoByteArray []byte) {
	leilao := common.ByteArrayToLeilao(leilaoByteArray)
	activeLeilao, ok := activeLeiloes[leilao.ID]
	if ok {
		lastLance := activeLeilao.LastValidLance
		if lastLance != (common.Lance{}) {
			q, err := common.CreateOrGetQueueAndBind("leilao_vencedor", chIn)
			common.FailOnError(err, "Error connecting to queue")

			common.PublishInQueue(chOut, q, common.LanceToByteArray(lastLance))
		}

		delete(activeLeiloes, leilao.ID)
	}
}

func consumeLeiloesFinalizados(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LEILAO FINALIZADO: %s", d.Body)

		go handleLeilaoFinalizado(d.Body)

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

	activeLeiloes = make(map[string]common.ActiveLeilao)

	qLance, err := common.CreateOrGetQueueAndBind("lance_realizado", chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLance, chIn, consomeLances)

	qLeiloesIniciados, err := common.CreateOrGetQueueAndBind("leilao_iniciado", chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesIniciados, chIn, consumeLeiloesIniciados)

	qLeiloesFinalizados, err := common.CreateOrGetQueueAndBind("leilao_finalizado", chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesFinalizados, chIn, consumeLeiloesFinalizados)

	var forever chan struct{}
	<-forever
}
