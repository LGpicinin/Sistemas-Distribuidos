package common

import (
	"crypto/sha256"
	"log"
)

func HashLance(lance Lance) []byte {
	lanceBytes := LanceToByteArray(lance)
	hash := sha256.New()
	_, err := hash.Write(lanceBytes)
	if err != nil {
		log.Fatalf("Error hashing message: %v", err)
	}

	hashedMessage := hash.Sum(nil)
	return hashedMessage
}
