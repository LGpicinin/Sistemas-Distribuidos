package main

import (
	common "common"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rabbitmq/amqp091-go"
)

var leiloes []common.Leilao
var leiloesSortedByStart []common.Leilao
var leiloesSortedByEnd []common.Leilao

func createLeiloes() []common.Leilao {
	ids := []string{"1", "2", "3"}
	bens := []string{"Geladeira", "Fusca", "Guarda Roupa"}

	for i := 0; i < 3; i++ {
		start := time.Now().Add(time.Second * time.Duration(rand.Int()%10*i))
		end := time.Now().Add(time.Second*time.Duration(rand.Int()%120) + time.Second*time.Duration(60))

		leilao := common.CreateLeilao(ids[i], bens[i], start, end)
		leiloes = append(leiloes, leilao)
	}

	return leiloes
}

func createFileLeiloes(leiloes []common.Leilao) {
	for i := range leiloes {
		leilao := leiloes[i]

		leilaoByteArray := common.LeilaoToByteArray(leilao)

		file, err := os.Create(fmt.Sprintf("ms-leilao/data/leilao-%s.json", leilao.ID))
		if err != nil {
			common.FailOnError(err, "Erro ao criar arquivo")
		}
		defer file.Close()

		file.Write(leilaoByteArray)
	}

}

func publishWhenStarts(ch *amqp091.Channel, q amqp091.Queue, leiloes []common.Leilao, allPublished chan bool) {
	for first := leiloes[0]; ; first = leiloes[0] {
		if time.Now().Compare(first.StartDate) >= 0 {
			common.PublishInQueue(ch, q, common.LeilaoToByteArray(first), "leilao_iniciado")
			leiloes = append(leiloes[:0], leiloes[1:]...)

			log.Printf("[MS-LEILAO] Published %s on %s\n\n", spew.Sdump(first), q.Name)
			common.CreateQueue(ch)
			if len(leiloes) == 0 {
				break
			}
		}
	}

	allPublished <- true
}

func publishWhenFinishes(ch *amqp091.Channel, q amqp091.Queue, leiloes []common.Leilao, allPublished chan bool) {
	for first := leiloes[0]; ; first = leiloes[0] {
		if time.Now().Compare(first.EndDate) >= 0 {
			common.PublishInQueue(ch, q, common.LeilaoToByteArray(first), "leilao_finalizado")
			leiloes = append(leiloes[:0], leiloes[1:]...)

			log.Printf("[MS-LEILAO] Published %s on %s", spew.Sdump(first), q.Name)
			if len(leiloes) == 0 {
				break
			}
		}
	}

	allPublished <- true
}

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	qIniciado, err := common.CreateOrGetQueueAndBind("leilao_iniciado", ch)
	common.FailOnError(err, "Error connecting to queue")
	qFinalizado, err := common.CreateOrGetQueueAndBind("leilao_finalizado", ch)
	common.FailOnError(err, "Error connecting to queue")

	leiloes = createLeiloes()

	createFileLeiloes(leiloes)

	leiloesSortedByStart = append(leiloesSortedByStart, leiloes...)
	sort.Sort(common.ByStartDate(leiloesSortedByStart))

	leiloesSortedByEnd = append(leiloesSortedByEnd, leiloes...)
	sort.Sort(common.ByEndDate(leiloesSortedByEnd))

	publishedAllLeiloesStart := make(chan bool)
	go publishWhenStarts(ch, qIniciado, leiloesSortedByStart, publishedAllLeiloesStart)

	publishedAllLeiloesEnd := make(chan bool)
	go publishWhenFinishes(ch, qFinalizado, leiloesSortedByEnd, publishedAllLeiloesEnd)

	<-publishedAllLeiloesStart
	<-publishedAllLeiloesEnd
}
