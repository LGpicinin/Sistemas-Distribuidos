package common

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"log"
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
		FailOnError(err, "Erro ao converter lance para []byte")
	}

	return leilaoByteArray
}

func (lance *Lance) FromByteArray(byteArray []byte) {
	err := json.Unmarshal(byteArray, lance)
	FailOnError(err, "Erro ao converter []byte para lance")
}

func (lance *Lance) Hash() []byte {
	lanceBytes := lance.ToByteArray()
	hash := sha256.New()
	_, err := hash.Write(lanceBytes)
	if err != nil {
		log.Fatalf("Error hashing message: %v", err)
	}

	hashedMessage := hash.Sum(nil)
	return hashedMessage
}

func HashLance(lance Lance) []byte {
	lanceBytes := lance.ToByteArray()
	hash := sha256.New()
	_, err := hash.Write(lanceBytes)
	if err != nil {
		log.Fatalf("Error hashing message: %v", err)
	}

	hashedMessage := hash.Sum(nil)
	return hashedMessage
}

func (lance *Lance) Sign(privateKey *rsa.PrivateKey) []byte {
	hashedLance := lance.Hash()
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashedLance)
	if err != nil {
		log.Fatalf("Error signing message: %v", err)
	}

	return signature
}
