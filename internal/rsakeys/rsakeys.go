package rsakeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

func GenerateRSAKeys() (privateKey []byte, publicKey []byte, err error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	privateKey = privateKeyToBytes(pk)
	publicKey = publicKeyToBytes(&pk.PublicKey)

	return
}

func privateKeyToBytes(privateKey *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(privateKey)
}

func publicKeyToBytes(publicKey *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(publicKey)
}