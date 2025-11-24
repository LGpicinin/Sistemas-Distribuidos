package main

import (
	common "common"
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	// dto "ms-leilao/DTO"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var activeLeiloes map[string]common.Leilao = make(map[string]common.Leilao)
var leiloesSortedByStart *list.List = list.New()
var leiloesSortedByEnd *list.List = list.New()

type createLeilaoHandler struct{}
type listLeilaoHandler struct{}

// função que insere novo leilão em lista ordenada por tempo
func insertionSortOnList(leilaoList *list.List, value common.Leilao, fieldToCompare string) {

	r := reflect.ValueOf(value)
	fieldValue := reflect.Indirect(r).FieldByName(fieldToCompare).Interface().(time.Time)

	if leilaoList.Len() == 0 {
		leilaoList.PushBack(value)
	} else {
		var k *list.Element = nil
		for e := leilaoList.Front(); e != nil; e = e.Next() {
			e_r := reflect.ValueOf(e.Value.(common.Leilao))
			e_value := reflect.Indirect(e_r).FieldByName(fieldToCompare).Interface().(time.Time)

			if e_value.Compare(fieldValue) > 0 {
				k = e
				break
			}
		}
		if k != nil {
			leilaoList.InsertBefore(value, k)
		} else {
			leilaoList.PushBack(value)
		}
	}
}

// função infinita que publica leilão na fila quando ele for iniciado
func publishWhenStarts(ch *amqp091.Channel, q amqp091.Queue) {
	for {
		if leiloesSortedByStart.Len() == 0 {
			continue
		}

		first := leiloesSortedByStart.Front()
		firstLeilao := first.Value.(common.Leilao)

		if !firstLeilao.HasStarted() {
			continue
		}

		common.PublishInQueue(ch, q, firstLeilao.ToByteArray(), common.QUEUE_LEILAO_INICIADO)
		activeLeiloes[firstLeilao.ID] = firstLeilao
		leiloesSortedByStart.Remove(first)

		log.Printf("[MS-LEILAO] NOVO LEILÃO INICIADO: %s PUBLICADO NA FILA %s\n\n", firstLeilao.Print(), q.Name)
	}
}

// função infinita que publica leilão na fila quando ele for finalizado
func publishWhenFinishes(ch *amqp091.Channel, q amqp091.Queue) {
	for {
		if leiloesSortedByEnd.Len() == 0 {
			continue
		}

		first := leiloesSortedByEnd.Front()
		firstLeilao := first.Value.(common.Leilao)

		if !firstLeilao.HasEnded() {
			continue
		}

		common.PublishInQueue(ch, q, firstLeilao.ToByteArray(), common.QUEUE_LEILAO_FINALIZADO)
		leiloesSortedByEnd.Remove(first)
		delete(activeLeiloes, firstLeilao.ID)

		log.Printf("[MS-LEILAO] NOVO LEILÃO FINALIZADO %s PUBLICADO NA FILA %s\n\n", firstLeilao.Print(), q.Name)
	}
}

// recebe requisição http do gateway para criação de novo leilão
// chama função para inserir na lista de ordenada por tempo de início e na de tempo por fim
func (h *createLeilaoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var leilao common.Leilao

	err := json.NewDecoder(r.Body).Decode(&leilao)
	if err != nil {
		log.Printf("Erro ao decodificar requisição: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	insertionSortOnList(leiloesSortedByStart, leilao, "StartDate")
	insertionSortOnList(leiloesSortedByEnd, leilao, "EndDate")

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(leilao)
}

// recebe requisição http do gateway para listar leilões ativos
// envia resposta por http
func (h *listLeilaoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var leiloes []common.Leilao
	for _, activeLeilao := range activeLeiloes {
		leiloes = append(leiloes, activeLeilao)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(leiloes)
}

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	qIniciado, err := common.CreateOrGetQueueAndBind("", common.QUEUE_LEILAO_INICIADO, ch)
	common.FailOnError(err, "Error connecting to queue")
	qFinalizado, err := common.CreateOrGetQueueAndBind(common.QUEUE_LEILAO_FINALIZADO, common.QUEUE_LEILAO_FINALIZADO, ch)
	common.FailOnError(err, "Error connecting to queue")

	go publishWhenStarts(ch, qIniciado)
	go publishWhenFinishes(ch, qFinalizado)

	mux := http.NewServeMux()
	mux.Handle("/create", &createLeilaoHandler{})
	mux.Handle("/list", &listLeilaoHandler{})
	fmt.Println("Server running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}
