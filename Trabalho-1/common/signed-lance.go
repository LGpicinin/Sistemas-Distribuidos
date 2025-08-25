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

func SignedLanceToByteArray(lance SignedLance) []byte {
	leilaoByteArray, err := json.Marshal(lance)
	if err != nil {
		FailOnError(err, "Erro ao converter lance para []byte")
	}

	return leilaoByteArray
}

func ByteArrayToSignedLance(byteArray []byte) SignedLance {
	var lance SignedLance
	err := json.Unmarshal(byteArray, &lance)
	if err != nil {
		FailOnError(err, "Erro ao converter []byte para lance")
	}

	return lance
}
