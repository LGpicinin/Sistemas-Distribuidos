package main

import (
	"fmt"
	"log"
	"os"

	common "common"

	"github.com/davecgh/go-spew/spew"
	"github.com/rabbitmq/amqp091-go"
)

func consomeLeilaoIniciado(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LEILAO] NOVO LEILÃO: %s", spew.Sdump(common.ByteArrayToLeilao(d.Body)))
		d.Ack(true)
	}
}

func hello() {
	fmt.Println("========== Bem vindo ao UTFPR Leilões ==========")
	fmt.Println("Caso deseje registrar um lance, pressione Enter. Caso deseje sair, aperte CTRL+C")
}

func publishLance(q amqp091.Queue, ch *amqp091.Channel, leilaoId string, userId string) {
	var value float32
	fmt.Print("Qual o valor do lance que deseja fazer? ")
	fmt.Scanf("%f", &value)

	lance := common.CreateLance(leilaoId, userId, value)

}

func menu(userId string, q amqp091.Queue, ch *amqp091.Channel) {
	var input string
	for fmt.Scanf("%s", &input); ; fmt.Scanf("%s", &input) {
		var leilaoId string
		fmt.Print("Digite o ID do Leilão em que deseja registrar um lance, e o valor do lance: ")
		fmt.Scanf("%s", &leilaoId)

	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Uso correto: ./client id_do_cliente")
	}

	userId := os.Args[1]

	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	q, err := common.CreateOrGetQueueAndBind("leilao_iniciado", ch)
	common.FailOnError(err, "Error connecting to queue")

	hello()
	common.ConsumeEvents(q, ch, consomeLeilaoIniciado)
	go menu(userId, q, ch)

	var forever chan struct{}
	<-forever
}
