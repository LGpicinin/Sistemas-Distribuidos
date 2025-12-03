package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func ReadAndParseKey(keyPath string) (*rsa.PublicKey, *rsa.PrivateKey) {
	keyFile, _ := os.ReadFile(keyPath)
	keyBlock, _ := pem.Decode(keyFile)

	var publicKeyBlock, privateKeyBlock *pem.Block
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		privateKeyBlock = keyBlock
		privateKey, _ := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		return nil, privateKey
	case "RSA PUBLIC KEY":
		publicKeyBlock = keyBlock
		publicKey, _ := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
		return publicKey, nil
	}

	return nil, nil
}
