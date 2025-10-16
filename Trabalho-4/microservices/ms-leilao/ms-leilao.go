package main

import (
	common "common"
	"container/list"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var activeLeiloes *list.List = list.New()

// var leiloes []common.Leilao
var leiloesSortedByStart *list.List = list.New()

// var leiloesSortedByStart []common.Leilao
var leiloesSortedByEnd *list.List = list.New()

// var leiloesSortedByEnd []common.Leilao

type createHandler struct{}
type listHandler struct{}

func publishWhenStarts(ch *amqp091.Channel, q amqp091.Queue, leiloes *list.List) {
	for {
		if leiloes.Len() > 0 {
			for first := leiloes.Front(); ; first = leiloes.Front() {
				firstLeilao := first.Value.(common.Leilao)
				if time.Now().Compare(firstLeilao.StartDate) >= 0 {
					common.PublishInQueue(ch, q, firstLeilao.ToByteArray(), common.QUEUE_LEILAO_INICIADO)
					activeLeiloes.PushBack(first)
					leiloes.Remove(first)

					log.Printf("[MS-LEILAO] NOVO LEILÃO INICIADO: %s PUBLICADO NA FILA %s\n\n", firstLeilao.Print(), q.Name)
					common.CreateQueue("", ch)
					if leiloes.Len() == 0 {
						break
					}
				}
			}
		}
	}

}

func publishWhenFinishes(ch *amqp091.Channel, q amqp091.Queue, leiloes *list.List) {
	for {
		if leiloes.Len() > 0 {
			for first := leiloes.Front(); ; first = leiloes.Front() {
				firstLeilao := first.Value.(common.Leilao)
				if time.Now().Compare(firstLeilao.EndDate) >= 0 {
					common.PublishInQueue(ch, q, firstLeilao.ToByteArray(), common.QUEUE_LEILAO_FINALIZADO)
					leiloes.Remove(first)

					log.Printf("[MS-LEILAO] NOVO LEILÃO FINALIZADO %s PUBLICADO NA FILA %s\n\n", firstLeilao.Print(), q.Name)
					if leiloes.Len() == 0 {
						break
					}
				}
			}
		}
	}
}

func (h *createHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var leilao common.Leilao

	err := json.NewDecoder(r.Body).Decode(&leilao)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if leiloesSortedByStart.Len() == 0 {
		leiloesSortedByStart.PushBack(leilao)
	} else {
		var k *list.Element
		for e := leiloesSortedByStart.Front(); e != nil; e = e.Next() {
			if e.Value.(common.Leilao).StartDate.Compare(leilao.StartDate) > 0 {
				k = e
				break
			}
		}
		leiloesSortedByStart.InsertBefore(leilao, k)
	}

	if leiloesSortedByEnd.Len() == 0 {
		leiloesSortedByEnd.PushBack(leilao)
	} else {
		var k *list.Element
		for e := leiloesSortedByEnd.Front(); e != nil; e = e.Next() {
			if e.Value.(common.Leilao).StartDate.Compare(leilao.StartDate) < 0 {
				k = e
				break
			}
		}
		leiloesSortedByEnd.InsertBefore(leilao, k)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(leilao.ToByteArray())
}

func (h *listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	qIniciado, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_INICIADO, ch)
	common.FailOnError(err, "Error connecting to queue")
	qFinalizado, err := common.CreateOrGetQueueAndBind(common.QUEUE_LEILAO_FINALIZADO, common.QUEUE_LEILAO_FINALIZADO, ch)
	common.FailOnError(err, "Error connecting to queue")

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/create", &createHandler{})
	mux.Handle("/list", &listHandler{})

	go publishWhenStarts(ch, qIniciado, leiloesSortedByStart)

	go publishWhenFinishes(ch, qFinalizado, leiloesSortedByEnd)

	http.ListenAndServe(":8090", mux)

}
