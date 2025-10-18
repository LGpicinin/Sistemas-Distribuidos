package main

import (
	common "common"
	"fmt"
	"log"

	"net/http"

	"github.com/rabbitmq/amqp091-go"
)

var connIn *amqp091.Connection
var chIn *amqp091.Channel

var connOut *amqp091.Connection
var chOut *amqp091.Channel

var activeLeiloes map[string]common.ActiveLeilao = make(map[string]common.ActiveLeilao)

type newLanceHandler struct{}

func handleLanceCandidate(lanceCanditate []byte) {
	var lance common.Lance
	lance.FromByteArray(lanceCanditate)

	activeLeilao, ok := activeLeiloes[lance.LeilaoID]
	if !ok {
		log.Printf("[MS-LANCE] Erro ao acessar leilão ativo: %v\n", lance.LeilaoID)
		return
	}

	if lance.Value <= activeLeilao.LastValidLance.Value {
		log.Printf("[MS-LANCE] Lance não válido: \n%s\n", lance.Print())

		q, err := common.CreateOrGetQueueAndBind(
			common.QUEUE_LANCE_INVALIDADO, common.QUEUE_LANCE_INVALIDADO, chIn,
		)
		common.FailOnError(err, "Error connecting to queue")
		common.PublishInQueue(chOut, q, lance.ToByteArray(), common.QUEUE_LANCE_INVALIDADO)
		return
	}

	activeLeilao.LastValidLance = lance

	activeLeiloes[lance.LeilaoID] = activeLeilao

	q, err := common.CreateOrGetQueueAndBind(common.QUEUE_LANCE_VALIDADO, common.QUEUE_LANCE_VALIDADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.PublishInQueue(chOut, q, lance.ToByteArray(), common.QUEUE_LANCE_VALIDADO)

	log.Printf("Novo lance validado: \n%s\n", lance.Print())
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
		log.Printf("[MS-LANCE] LEILÃO FINALIZADO: \n%s\n", leilao.Print())
	}
}

func consumeLeiloesFinalizados(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		// log.Printf("[MS-LANCE] NOVO LEILAO FINALIZADO: %s", d.Body)

		go handleLeilaoFinalizado(d.Body)

		d.Ack(false)
	}
}

func (h *newLanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	mux.Handle("/create", &newLanceHandler{})

	qLance, err := common.CreateOrGetQueueAndBind(common.QUEUE_LANCE_REALIZADO, common.QUEUE_LANCE_REALIZADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLance, chIn, consomeLances)

	qLeiloesIniciados, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_INICIADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesIniciados, chIn, consumeLeiloesIniciados)

	qLeiloesFinalizados, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_FINALIZADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesFinalizados, chIn, consumeLeiloesFinalizados)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
