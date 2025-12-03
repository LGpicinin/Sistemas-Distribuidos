package common

import (
	"encoding/json"
	"fmt"
)

type Payment struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
	UserID   string  `json:"clientId"`
	Adress   string  `json:"webhookAddress"`
}

func CreatePayment(currency string, userId string, value float32, webhookAddress string) Payment {
	var payment Payment = Payment{
		Value:    value,
		Currency: currency,
		UserID:   userId,
		Adress:   webhookAddress,
	}

	return payment
}

func (payment *Payment) ToByteArray() []byte {
	paymentByteArray, err := json.Marshal(*payment)
	if err != nil {
		FailOnError(err, "Erro ao converter payment para []byte")
	}

	return paymentByteArray
}

func (payment *Payment) FromByteArray(byteArray []byte) {
	err := json.Unmarshal(byteArray, payment)
	FailOnError(err, "Erro ao converter []byte para payment")
}

func (payment *Payment) Print() string {
	return "Pagamento:\n" +
		"\tID do usuário: " + payment.UserID + "\n" +
		"\tValor do lance: R$ " + fmt.Sprintf("%f", payment.Value) + "\n" +
		"\tMoeda: " + payment.Currency + "\n"
}
