package models

import (
	common "common/utils"
	"encoding/json"
	"fmt"
)

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

func (lance *Lance) ToByteArray() []byte {
	leilaoByteArray, err := json.Marshal(*lance)
	if err != nil {
		common.FailOnError(err, "Erro ao converter lance para []byte")
	}

	return leilaoByteArray
}

func (lance *Lance) FromByteArray(byteArray []byte) {
	err := json.Unmarshal(byteArray, lance)
	common.FailOnError(err, "Erro ao converter []byte para lance")
}

func (lance *Lance) Print() string {
	return "Lance:\n" +
		"\tID do Leilão: " + lance.LeilaoID + "\n" +
		"\tID do usuário: " + lance.UserID + "\n" +
		"\tValor do Lance: R$ " + fmt.Sprintf("%f", lance.Value) + "\n"
}
