package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"

	common "common"

	"github.com/davecgh/go-spew/spew"
	"github.com/rabbitmq/amqp091-go"
)

var leiloesInteressados map[string]string

func consomeLeilaoIniciado(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		log.Printf("[MS-LEILAO] NOVO LEILÃO: %s", spew.Sdump(common.ByteArrayToLeilao(d.Body)))
		d.Ack(false)
	}
}

func hello() {
	fmt.Println("========== Bem vindo ao UTFPR Leilões ==========")
	fmt.Println("Caso deseje registrar um lance, pressione Enter. Caso deseje sair, aperte CTRL+C")
}

func handleNotificacao(notificacaoByteArray []byte) {
	notificacao := common.ByteArrayToNotificacao((notificacaoByteArray))

	if notificacao.Status == common.NovoLance {
		log.Printf("[MS-LEILAO] NOVA NOTIFICAÇÃO LANCE: %s", spew.Sdump(notificacao))
	} else {
		delete(leiloesInteressados, notificacao.Lance.LeilaoID)

		log.Printf("[MS-LEILAO] NOVA NOTIFICAÇÃO GANHADOR: %s", spew.Sdump(notificacao))
	}

}

func consomeLeilaoInteressado(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		go handleNotificacao(d.Body)

		d.Ack(false)
	}
}

func publishLance(q amqp091.Queue, ch *amqp091.Channel, leilaoId string, userId string, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) {
	var value float32
	fmt.Print("Qual o valor do lance que deseja fazer? ")
	fmt.Scanf("%f", &value)

	lance := common.CreateLance(leilaoId, userId, value)
	signature := signLance(lance, privateKey)

	signedLance := common.CreateSignedLance(lance, signature)
	signedLanceBytes := common.SignedLanceToByteArray(signedLance)

	common.PublishInQueue(ch, q, signedLanceBytes, "lance_realizado")

	_, ok := leiloesInteressados[leilaoId]
	if !ok {
		leiloesInteressados[leilaoId] = leilaoId

		nome_fila := fmt.Sprintf("leilao_%s", leilaoId)
		q, err := common.CreateOrGetQueueAndBind(nome_fila, ch)
		common.FailOnError(err, "Error connecting to queue")

		common.ConsumeEvents(q, ch, consomeLeilaoInteressado)
	}
}

func signLance(lance common.Lance, privateKey *rsa.PrivateKey) []byte {
	hashedLance := common.HashLance(lance)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashedLance)
	if err != nil {
		log.Fatalf("Error signing message: %v", err)
	}

	return signature
}

func menu(userId string, q amqp091.Queue, ch *amqp091.Channel, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) {
	var input string
	for fmt.Scanf("%s", &input); ; fmt.Scanf("%s", &input) {
		var leilaoId string
		fmt.Print("Digite o ID do Leilão em que deseja registrar um lance: ")
		fmt.Scanf("%s", &leilaoId)

		publishLance(q, ch, leilaoId, userId, publicKey, privateKey)
	}
}

func createKeys(userId string) (*rsa.PublicKey, *rsa.PrivateKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Error generating key pair: %v", err)
	}
	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile(fmt.Sprintf("./client/keys/private/%s.pem", userId), privateKeyPEM, 0600)
	if err != nil {
		log.Fatalf("Error saving private key: %v", err)
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile(fmt.Sprintf("./client/keys/public/%s.pem", userId), publicKeyPEM, 0644)
	if err != nil {
		log.Fatalf("Error saving public key: %v", err)
	}

	err = os.WriteFile(fmt.Sprintf("./ms-lance/keys/public/%s.pem", userId), publicKeyPEM, 0644)
	if err != nil {
		log.Fatalf("Error saving public key: %v", err)
	}

	return publicKey, privateKey
}

func initClient(args []string) (*rsa.PublicKey, *rsa.PrivateKey, string) {
	userId := args[1]

	entries, _ := os.ReadDir("./client/keys/public")
	for _, entry := range entries {
		if strings.Contains(entry.Name(), userId) {
			publicKey, _ := common.ReadAndParseKey(fmt.Sprintf("%s/%s", "./client/keys/public", entry.Name()))
			_, privateKey := common.ReadAndParseKey(fmt.Sprintf("%s/%s", "./client/keys/private", entry.Name()))

			return publicKey, privateKey, userId
		}
	}

	publicKey, privateKey := createKeys(userId)
	return publicKey, privateKey, userId
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Uso correto: ./bin/client <entre|crie_id_do_cliente>")
	}

	publicKey, privateKey, userId := initClient(os.Args)
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	leiloesInteressados = make(map[string]string)

	q, err := common.CreateOrGetQueueAndBind("leilao_iniciado", ch)
	common.FailOnError(err, "Error connecting to queue")

	hello()
	common.ConsumeEvents(q, ch, consomeLeilaoIniciado)

	qLances, err := common.CreateOrGetQueueAndBind("lance_realizado", ch)
	common.FailOnError(err, "Error connecting to queue")
	go menu(userId, qLances, ch, publicKey, privateKey)

	var forever chan struct{}
	<-forever
}
