package models

import (
	common "common/utils"
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

func (notificacao *Notificacao) ToByteArray() []byte {
	notificacaoByteArray, err := json.Marshal(*notificacao)
	common.FailOnError(err, "Erro ao converter notificação para []byte")

	return notificacaoByteArray
}

func (notificacao *Notificacao) FromByteArray(byteArray []byte) {
	err := json.Unmarshal(byteArray, notificacao)
	common.FailOnError(err, "Erro ao converter []byte para notificação")
}

func (notificacao *Notificacao) Print() string {
	return "\nNotificação:\n" +
		"\tStatus: " + string(notificacao.Status) + "\n" +
		"\t" + notificacao.Lance.Print() + "\n"
}
