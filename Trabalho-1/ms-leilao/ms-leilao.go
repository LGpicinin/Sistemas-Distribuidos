package main

import (
	common "common"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

var leiloes []common.Leilao

func createLeiloes() []common.Leilao {

	start := time.Now().Add(10000)
	end := time.Now().Add(30000)

	ids := []string{"1", "2"}
	bens := []string{"Geladeira", "Fusca"}

	var i int

	for i = 0; i < 2; i++ {
		leilao := common.CreateLeilao(ids[i], bens[i], start, end)
		leiloes = append(leiloes, leilao)
	}

	return leiloes
}

func createFileLeiloes(leiloes []common.Leilao) {
	var data [][]string

	data = common.LeiloesToCsv(leiloes)

	// create a file
	file, err := os.Create("ms-leilao/data/leiloes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)

	defer writer.Flush()

	// write all rows at once
	writer.WriteAll(data)

	file.Close()

}

func readLeiloesFile() []byte {
	content, err := os.ReadFile("ms-leilao/data/leiloes.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))

	return content
}

func main() {
	conn, ch := common.ConnectToBroker()
	defer conn.Close()
	defer ch.Close()

	qIniciado, err := common.CreateOrGetQueueAndBind("leilao_iniciado", ch)
	common.FailOnError(err, "Error connecting to queue")

	leiloes = createLeiloes()

	createFileLeiloes(leiloes)

	body := readLeiloesFile()

	common.PublishInQueue(ch, qIniciado, body)

	log.Printf(" [x] Sent %s", body)
}
