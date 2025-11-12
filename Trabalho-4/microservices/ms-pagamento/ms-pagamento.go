package main

import (
	"bytes"
	common "common"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var connIn *amqp091.Connection
var chIn *amqp091.Channel

var connOut *amqp091.Connection
var chOut *amqp091.Channel

type statusPagamentoHandler struct{}

func handleLeilaoGanhador(lanceByteArray []byte) {
	var lance common.Lance
	lance.FromByteArray(lanceByteArray)

	var payment common.Payment = common.CreatePayment("R$", lance.UserID, lance.Value, "http://localhost:8100/status")

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3333/create-payment", bytes.NewReader(payment.ToByteArray()))

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("request failed: %s", err)
	}
	defer res.Body.Close()

	log.Printf("[MS-PAGAMENTO] Novo pagamento gerado: \n%s\n", payment.Print())

	var link common.JustLink

	err = json.NewDecoder(res.Body).Decode(&link)

	var link_data common.Link = common.CreateLink(link.Link, lance.UserID)

	q, err := common.CreateOrGetQueueAndBind(
		common.QUEUE_LINK_PAGAMENTO, common.QUEUE_LINK_PAGAMENTO, chIn,
	)
	common.FailOnError(err, "Error connecting to queue")
	common.PublishInQueue(chOut, q, link_data.ToByteArray(), common.QUEUE_LINK_PAGAMENTO)

	log.Printf("[MS-PAGAMENTO] Novo link gerado: \n%s\n", link_data.Print())

}

func consumeLeiloesGanhador(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		// log.Printf("[MS-PAGAMENTO] NOVO LEILAO INICIADO: %s", d.Body)

		go handleLeilaoGanhador(d.Body)

		d.Ack(false)
	}
}

func (h *statusPagamentoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var status common.StatusPayment

	err := json.NewDecoder(r.Body).Decode(&status)

	q, err := common.CreateOrGetQueueAndBind(
		common.QUEUE_STATUS_PAGAMENTO, common.QUEUE_STATUS_PAGAMENTO, chIn,
	)
	common.FailOnError(err, "Error connecting to queue")
	common.PublishInQueue(chOut, q, status.ToByteArray(), common.QUEUE_STATUS_PAGAMENTO)

	log.Printf("[MS-PAGAMENTO] Novo pagamento gerado: \n%s\n", status.Print())
}

func main() {
	connIn, chIn = common.ConnectToBroker()
	defer connIn.Close()
	defer chIn.Close()

	connOut, chOut = common.ConnectToBroker()
	defer connOut.Close()
	defer chOut.Close()

	qLeiloesGanhador, err := common.CreateOrGetQueueAndBind(
		"", common.QUEUE_LEILAO_VENCEDOR, chIn,
	) // necessário verificar se pode ou não nomear fila
	common.FailOnError(err, "Error connecting to queue")
	common.ConsumeEvents(qLeiloesGanhador, chIn, consumeLeiloesGanhador)

	q1, err := common.CreateOrGetQueueAndBind(
		common.QUEUE_LINK_PAGAMENTO, common.QUEUE_LINK_PAGAMENTO, chIn,
	)
	q2, err := common.CreateOrGetQueueAndBind(
		common.QUEUE_STATUS_PAGAMENTO, common.QUEUE_STATUS_PAGAMENTO, chIn,
	)
	fmt.Println(q1.Name)
	fmt.Println(q2.Name)

	mux := http.NewServeMux()
	mux.Handle("/status", &statusPagamentoHandler{})
	fmt.Println("Server running on http://localhost:8100")
	http.ListenAndServe(":8100", mux)
}
