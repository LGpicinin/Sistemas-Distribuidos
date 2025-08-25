package common

import "encoding/json"

type Lance struct {
	LeilaoID string  `json:"leilao_id"`
	UserID   string  `json:"user_id"`
	Value    float32 `json:"value"`
}

func CreateLance(leilaoId string, userId string, value float32) Lance {
	var lance Lance = Lance{
		LeilaoID: leilaoId,
		UserID:   userId,
		Value:    value,
	}

	return lance
}

func LanceToByteArray(lance Lance) []byte {
	leilaoByteArray, err := json.Marshal(lance)
	if err != nil {
		FailOnError(err, "Erro ao converter lance para []byte")
	}

	return leilaoByteArray
}

func ByteArrayToLance(byteArray []byte) Lance {
	var lance Lance
	err := json.Unmarshal(byteArray, &lance)
	if err != nil {
		FailOnError(err, "Erro ao converter []byte para lance")
	}

	return lance
}
