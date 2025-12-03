package common

import (
	"encoding/json"
	"fmt"
)

type StatusPayment struct {
	Value     float32 `json:"value"`
	PaymentId string  `json:"paymentId"`
	UserID    string  `json:"clientId"`
	Status    bool    `json:"status"`
}

func (status *StatusPayment) ToByteArray() []byte {
	statusByteArray, err := json.Marshal(*status)
	if err != nil {
		FailOnError(err, "Erro ao converter payment para []byte")
	}

	return statusByteArray
}

func (status *StatusPayment) Print() string {
	return "Pagamento:\n" +
		"\tID do usuário: " + status.UserID + "\n" +
		"\tValor do lance: R$ " + fmt.Sprintf("%f", status.Value) + "\n" +
		"\tID do Pagamento: " + status.PaymentId + "\n" +
		"\tStatus: " + fmt.Sprintf("%t", status.Status) + "\n"
}
