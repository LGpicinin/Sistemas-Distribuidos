package main

import (
	common "common"
	"encoding/json"
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

func handleLanceCandidate(lanceCandidate common.Lance) {
	activeLeilao, ok := activeLeiloes[lanceCandidate.LeilaoID]
	if !ok {
		log.Printf("[MS-LANCE] Erro ao acessar leilão ativo: %v\n", lanceCandidate.LeilaoID)
		return
	}

	if lanceCandidate.Value <= activeLeilao.LastValidLance.Value {
		log.Printf("[MS-LANCE] Lance não válido: \n%s\n", lanceCandidate.Print())

		q, err := common.CreateOrGetQueueAndBind(
			common.QUEUE_LANCE_INVALIDADO, common.QUEUE_LANCE_INVALIDADO, chIn,
		)
		common.FailOnError(err, "Error connecting to queue")
		common.PublishInQueue(
			chOut, q, lanceCandidate.ToByteArray(), common.QUEUE_LANCE_INVALIDADO,
		)
		return
	}

	activeLeilao.LastValidLance = lanceCandidate

	activeLeiloes[lanceCandidate.LeilaoID] = activeLeilao

	q, err := common.CreateOrGetQueueAndBind(
		common.QUEUE_LANCE_VALIDADO, common.QUEUE_LANCE_VALIDADO, chIn,
	)
	common.FailOnError(err, "Error connecting to queue")
	common.PublishInQueue(chOut, q, lanceCandidate.ToByteArray(), common.QUEUE_LANCE_VALIDADO)

	log.Printf("[MS-LANCE] Novo lance validado: \n%s\n", lanceCandidate.Print())
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
			q, err := common.CreateOrGetQueueAndBind(
				common.QUEUE_LEILAO_VENCEDOR, common.QUEUE_LEILAO_VENCEDOR, chOut,
			)
			common.FailOnError(err, "Error connecting to queue")

			common.PublishInQueue(chOut, q, lastLance.ToByteArray(), common.QUEUE_LEILAO_VENCEDOR)
			log.Printf("[MS-LANCE] NOVO VENCEDOR: \n%s\n", lastLance.Print())
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var lance common.Lance

	err := json.NewDecoder(r.Body).Decode(&lance)
	if err != nil {
		log.Printf("Erro ao decodificar requisição: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handleLanceCandidate(lance)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(lance)
}

func main() {
	connIn, chIn = common.ConnectToBroker()
	defer connIn.Close()
	defer chIn.Close()

	connOut, chOut = common.ConnectToBroker()
	defer connOut.Close()
	defer chOut.Close()

	qLeiloesIniciados, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_INICIADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesIniciados, chIn, consumeLeiloesIniciados)

	qLeiloesFinalizados, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_FINALIZADO, chIn)
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesFinalizados, chIn, consumeLeiloesFinalizados)

	_, _ = common.CreateOrGetQueueAndBind(
		common.QUEUE_LANCE_INVALIDADO, common.QUEUE_LANCE_INVALIDADO, chIn,
	)

	mux := http.NewServeMux()
	mux.Handle("/new", &newLanceHandler{})
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
