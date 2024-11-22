package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrUnableToDecodePrivateKeyPEMBlock = errors.New("failed to decode PEM block containing private key")
var ErrUnableToDecodePublicKeyPEMBlock = errors.New("failed to decode PEM block containing public key")

type ECDSAKeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// GenerateECDSAKeyPair generates an ECDSA private and public key pair
func New() (*ECDSAKeyPair, error) {
	// Use the P-256 curve (also known as prime256v1)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ECDSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// NewFromBytes creates an ECDSA key pair from the given private and public key bytes
func NewFromBytes(privateKeyBytes []byte) (*ECDSAKeyPair, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil || privateKeyBlock.Type != "EC PRIVATE KEY" {
		return nil, ErrUnableToDecodePrivateKeyPEMBlock
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECDSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// SavePrivateKey saves the ECDSA private key to a PEM file
func (k *ECDSAKeyPair) PrivateKeyPEM() ([]byte, error) {
	// Marshal the private key into ASN.1 DER format
	keyBytes, err := x509.MarshalECPrivateKey(k.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Create a PEM block for the private key
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})

	return privateKeyPEM, nil
}

// SavePublicKey saves the ECDSA public key to a PEM file
func (k *ECDSAKeyPair) PublicKeyPEM() ([]byte, error) {
	// Marshal the public key into ASN.1 DER format
	keyBytes, err := x509.MarshalPKIXPublicKey(k.PublicKey)
	if err != nil {
		return nil, err
	}

	// Create a PEM block for the public key
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	})

	return publicKeyPEM, nil
}
