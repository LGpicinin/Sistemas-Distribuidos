package common

import "encoding/json"

type SignedLance struct {
	Lance     Lance  `json:"lance"`
	Signature []byte `json:"signature"`
}

func CreateSignedLance(lance Lance, signature []byte) SignedLance {
	var signedLance SignedLance = SignedLance{
		Lance:     lance,
		Signature: signature,
	}

	return signedLance
}

func (signedLance *SignedLance) ToByteArray() []byte {
	leilaoByteArray, err := json.Marshal(*signedLance)
	FailOnError(err, "Erro ao converter lance para []byte")

	return leilaoByteArray
}

func (signedLance *SignedLance) FromByteArray(byteArray []byte) {
	err := json.Unmarshal(byteArray, signedLance)
	FailOnError(err, "Erro ao converter []byte para lance")
}
