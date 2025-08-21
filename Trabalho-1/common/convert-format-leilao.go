package common

import "encoding/json"

func LeilaoToByteArray(leilao Leilao) []byte {
	leilaoByteArray, err := json.Marshal(leilao)
	if err != nil {
		FailOnError(err, "Erro ao converter leilao para []byte")
	}

	return leilaoByteArray
}

func ByteArrayToLeilao(byteArray []byte) Leilao {
	var leilao Leilao
	err := json.Unmarshal(byteArray, &leilao)
	if err != nil {
		FailOnError(err, "Erro ao converter []byte para leilao")
	}

	return leilao
}
