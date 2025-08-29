package common

import (
	"encoding/json"
)

type Status string

const (
	NovoLance     Status = "Novo Lance"
	GanhadorLance Status = "Ganhador Lance"
)

type Notificacao struct {
	Status Status `json:"status"`
	Lance  Lance  `json:"lance"`
}

func CreateNotificacao(lance Lance, status Status) Notificacao {
	var notificacao Notificacao = Notificacao{
		Status: status,
		Lance:  lance,
	}

	return notificacao
}

func NotificacaoToByteArray(notificacao Notificacao) []byte {
	notificacaoByteArray, err := json.Marshal(notificacao)
	if err != nil {
		FailOnError(err, "Erro ao converter notificação para []byte")
	}

	return notificacaoByteArray
}

func ByteArrayToNotificacao(byteArray []byte) Notificacao {
	var notificacao Notificacao
	err := json.Unmarshal(byteArray, &notificacao)
	if err != nil {
		FailOnError(err, "Erro ao converter []byte para notificação")
	}

	return notificacao
}
