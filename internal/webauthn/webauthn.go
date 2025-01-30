package webauthn

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
)

type Credential struct {
	ID        []byte
	PublicKey any // *ecdsa.PublicKey or *rsa.PublicKey
	AAGUID    string
}

type ParseRequest struct {
	RPID              string
	AttestationObject string
}

func Parse(req *ParseRequest) (*Credential, error) {
	attestation, err := base64.RawURLEncoding.DecodeString(req.AttestationObject)
	if err != nil {
		return nil, err
	}

	var attestationData struct {
		AuthData []byte `cbor:"authData"`
	}
	if err := cbor.Unmarshal(attestation, &attestationData); err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(attestationData.AuthData)

	rpIDSHA256 := sha256.Sum256([]byte(req.RPID))
	rpHash := b.Next(32) // 32-byte rp hash

	if !bytes.Equal(rpHash, rpIDSHA256[:]) {
		return nil, fmt.Errorf("invalid rp id")
	}

	flags := b.Next(1) // flags
	if flags[0]&(0x1<<6) == 0 {
		return nil, fmt.Errorf("attestation object must have AT flag set")
	}

	_ = b.Next(4) // signature counter

	aaguidBytes := b.Next(16) // aaguid
	aaguid := uuid.UUID(aaguidBytes).String()

	lenCredID := binary.BigEndian.Uint16(b.Next(2)) // credential id len
	id := b.Next(int(lenCredID))                    // n bytes of credential id

	// remaining data is cbor-encoded COSE key
	var coseKey struct {
		Alg int    `cbor:"3,keyasint"`
		Crv int    `cbor:"-1,keyasint"`
		X   []byte `cbor:"-2,keyasint"`
		Y   []byte `cbor:"-3,keyasint"`
	}
	if err := cbor.NewDecoder(b).Decode(&coseKey); err != nil {
		return nil, err
	}

	var pub any
	switch coseKey.Alg {
	case -7:
		if coseKey.Crv != 1 {
			return nil, fmt.Errorf("unsupported curve")
		}

		pub = &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     big.NewInt(0).SetBytes(coseKey.X),
			Y:     big.NewInt(0).SetBytes(coseKey.Y),
		}
	default:
		return nil, fmt.Errorf("unsupported algorithm")
	}

	return &Credential{
		ID:        id,
		PublicKey: pub,
		AAGUID:    aaguid,
	}, nil
}

type VerifyRequest struct {
	RPID              string
	Origin            string
	ChallengeSHA256   []byte
	ClientDataJSON    string
	AuthenticatorData string
	Signature         string
}

func (c *Credential) Verify(req *VerifyRequest) error {
	clientDataBytes, err := base64.RawURLEncoding.DecodeString(req.ClientDataJSON)
	if err != nil {
		return err
	}

	clientDataHash := sha256.Sum256(clientDataBytes)

	authenticatorDataBytes, err := base64.RawURLEncoding.DecodeString(req.AuthenticatorData)
	if err != nil {
		return err
	}

	signatureBytes, err := base64.RawURLEncoding.DecodeString(req.Signature)
	if err != nil {
		return err
	}

	var signedBytes []byte
	signedBytes = append(signedBytes, authenticatorDataBytes...)
	signedBytes = append(signedBytes, clientDataHash[:]...)
	hash := sha256.Sum256(signedBytes)

	var sig struct {
		R, S *big.Int
	}
	if _, err := asn1.Unmarshal(signatureBytes, &sig); err != nil {
		return err
	}

	if !ecdsa.Verify(c.PublicKey.(*ecdsa.PublicKey), hash[:], sig.R, sig.S) {
		return fmt.Errorf("invalid signature")
	}

	// verify rp id hash
	b := bytes.NewBuffer(authenticatorDataBytes)
	rpIDSHA256 := sha256.Sum256([]byte(req.RPID))
	rpHash := b.Next(32) // 32-byte rp hash
	if !bytes.Equal(rpHash, rpIDSHA256[:]) {
		return fmt.Errorf("invalid rp id")
	}

	flags := b.Next(1) // flags
	if flags[0]&(0x1) == 0 {
		return fmt.Errorf("authenticator data must have UP flag set")
	}

	var clientData struct {
		Type        string `json:"type"`
		Challenge   string `json:"challenge"`
		Origin      string `json:"origin"`
		CrossOrigin bool   `json:"crossOrigin"`
	}
	if err := json.Unmarshal(clientDataBytes, &clientData); err != nil {
		return err
	}

	if clientData.Type != "webauthn.get" {
		return fmt.Errorf("invalid client data type")
	}

	if clientData.CrossOrigin {
		return fmt.Errorf("cross-origin not supported")
	}

	if clientData.Origin != req.Origin {
		return fmt.Errorf("invalid origin")
	}

	challengeBytes, err := base64.RawURLEncoding.DecodeString(clientData.Challenge)
	if err != nil {
		return err
	}

	challengeBytesSHA256 := sha256.Sum256(challengeBytes)
	if !bytes.Equal(challengeBytesSHA256[:], req.ChallengeSHA256) {
		return fmt.Errorf("invalid challenge")
	}

	return nil
}
