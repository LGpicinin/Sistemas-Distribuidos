package main

import (
	common "common"
	"crypto"
	"crypto/rsa"
	"fmt"
	"log"

	// "github.com/davecgh/go-spew/spew"
	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var ch *amqp091.Channel

// var activeLeiloes map[string]common.Leilao

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

	log.Printf("OIOI")
}

func consomeLances(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LANCE] NOVO LANCE: %s", d.Body)

		go handleLanceCandidate(d.Body)

		d.Ack(true)
	}
}

func main() {
	conn, ch = common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	// activeLeiloes = make(map[string]common.Leilao)

	qLance, err := common.CreateOrGetQueueAndBind("lance_realizado", ch)
	common.FailOnError(err, "Error connecting to queue")

	common.ConsumeEvents(qLance, ch, consomeLances)

	var forever chan struct{}
	<-forever
}
