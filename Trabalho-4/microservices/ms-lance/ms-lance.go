package main

import (
	common "common"
	"crypto"
	"crypto/rsa"
	"fmt"
	"log"

	"net/http"

	"github.com/rabbitmq/amqp091-go"
)

var connIn *amqp091.Connection
var chIn *amqp091.Channel

var connOut *amqp091.Connection
var chOut *amqp091.Channel

var activeLeiloes map[string]common.ActiveLeilao

type createHandler struct{}
type listHandler struct{}

func verifySignature(signedLance common.SignedLance) (bool, error) {
	hashedLance := signedLance.Lance.Hash()

	publicKey, _ := common.ReadAndParseKey(fmt.Sprintf("ms-lance/keys/public/%s.pem", signedLance.Lance.UserID))

	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedLance, signedLance.Signature)
	return err == nil, err
}

func handleLanceCandidate(lanceCanditate []byte) {
	var signedLance common.SignedLance
	signedLance.FromByteArray(lanceCanditate)

	isValidSignature, err := verifySignature(signedLance)
	if !isValidSignature {
		log.Printf("[MS-LANCE] Erro ao verificar chave: %v\n", err)
		return
	}

	activeLeilao, ok := activeLeiloes[signedLance.Lance.LeilaoID]
	if !ok {
		log.Printf("[MS-LANCE] Erro ao acessar leilão ativo: %v\n", signedLance.Lance.LeilaoID)
		return
	}

	if signedLance.Lance.Value <= activeLeilao.LastValidLance.Value {
		log.Printf("[MS-LANCE] Lance não válido: \n%s\n", signedLance.Lance.Print())
		return
	}

	activeLeilao.LastValidLance = signedLance.Lance

	activeLeiloes[signedLance.Lance.LeilaoID] = activeLeilao

	q, err := common.CreateOrGetQueueAndBind(common.QUEUE_LANCE_VALIDADO, common.QUEUE_LANCE_VALIDADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.PublishInQueue(chOut, q, signedLance.Lance.ToByteArray(), common.QUEUE_LANCE_VALIDADO)

	log.Printf("Novo lance validado: \n%s\n", signedLance.Lance.Print())
}

func consomeLances(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		// log.Printf("[MS-LANCE] NOVO LANCE: %s", d.Body)

		go handleLanceCandidate(d.Body)

		d.Ack(false)
	}
}

func handleLeilaoIniciado(leilaoByteArray []byte) {
	var leilao common.Leilao
	leilao.FromByteArray(leilaoByteArray)

	activeLeiloes[leilao.ID] = common.ActiveLeilao{
		Leilao: leilao,
	}

	log.Printf("[MS-LANCE] NOVO LEILÃO INICIADO: \n%s\n", leilao.Print())
}

func consumeLeiloesIniciados(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		// log.Printf("[MS-LANCE] NOVO LEILAO INICIADO: %s", d.Body)

		go handleLeilaoIniciado(d.Body)

		d.Ack(false)
	}
}

func handleLeilaoFinalizado(leilaoByteArray []byte) {
	var leilao common.Leilao
	leilao.FromByteArray(leilaoByteArray)

	activeLeilao, ok := activeLeiloes[leilao.ID]
	if ok {
		lastLance := activeLeilao.LastValidLance
		if lastLance != (common.Lance{}) {
			q, err := common.CreateOrGetQueueAndBind(common.QUEUE_LEILAO_VENCEDOR, common.QUEUE_LEILAO_VENCEDOR, chOut)
			common.FailOnError(err, "Error connecting to queue")

			common.PublishInQueue(chOut, q, lastLance.ToByteArray(), common.QUEUE_LEILAO_VENCEDOR)
		}

		delete(activeLeiloes, leilao.ID)
		log.Printf("[MS-LEILAO] LEILÃO FINALIZADO: \n%s\n", leilao.Print())
	}
}

func consumeLeiloesFinalizados(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		// log.Printf("[MS-LANCE] NOVO LEILAO FINALIZADO: %s", d.Body)

		go handleLeilaoFinalizado(d.Body)

		d.Ack(false)
	}
}

func (h *createHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

func (h *listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

func main() {
	connIn, chIn = common.ConnectToBroker()
	defer connIn.Close()
	defer chIn.Close()

	connOut, chOut = common.ConnectToBroker()
	defer connOut.Close()
	defer chOut.Close()

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/create", &createHandler{})
	mux.Handle("/list", &listHandler{})

	activeLeiloes = make(map[string]common.ActiveLeilao)

	qLance, err := common.CreateOrGetQueueAndBind(common.QUEUE_LANCE_REALIZADO, common.QUEUE_LANCE_REALIZADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLance, chIn, consomeLances)

	qLeiloesIniciados, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_INICIADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesIniciados, chIn, consumeLeiloesIniciados)

	qLeiloesFinalizados, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_FINALIZADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesFinalizados, chIn, consumeLeiloesFinalizados)

	http.ListenAndServe(":8080", mux)

	var forever chan struct{}
	<-forever
}
