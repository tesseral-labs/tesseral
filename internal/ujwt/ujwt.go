// Package ujwt implements a micro-subset of JSON Web Tokens.
//
// Only unencrypted, ES256-signed JWTs are supported. Only the `kid` header can
// be extracted. Claims cannot be extracted without verification.
package ujwt

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

var ErrBadJWT = errors.New("invalid jwt")

type headerData struct {
	KID string `json:"kid"`
	Alg string `json:"alg"`
}

type stdClaims struct {
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Nbf int64  `json:"nbf"`
}

// KeyID returns the `kid` header in token. This header value may be spoofed.
// KeyID does not authenticate anything.
//
// If the token is malformed, KeyID returns an error.
func KeyID(token string) (string, error) {
	headerBytes, _, _, err := parse(token)
	if err != nil {
		return "", err
	}

	var header headerData
	if err := jsonUnmarshalBase64(headerBytes, &header); err != nil {
		return "", ErrBadJWT
	}

	if header.Alg != "ES256" {
		return "", ErrBadJWT
	}

	return header.KID, nil
}

func Claims(pub *ecdsa.PublicKey, aud string, now time.Time, v any, token string) error {
	headerBytes, claimsBytes, signatureBytes, err := parse(token)
	if err != nil {
		return err
	}

	var header headerData
	if err := jsonUnmarshalBase64(headerBytes, &header); err != nil {
		return ErrBadJWT
	}

	if header.Alg != "ES256" {
		return ErrBadJWT
	}

	var signedPart []byte
	signedPart = append(signedPart, headerBytes...)
	signedPart = append(signedPart, '.')
	signedPart = append(signedPart, claimsBytes...)

	hash := sha256.Sum256(signedPart)

	sig, err := base64.RawURLEncoding.DecodeString(string(signatureBytes))
	if err != nil {
		return ErrBadJWT
	}

	if len(sig) != 64 {
		return ErrBadJWT
	}

	var r, s big.Int
	r.SetBytes(sig[:32])
	s.SetBytes(sig[32:])

	if !ecdsa.Verify(pub, hash[:], &r, &s) {
		return ErrBadJWT
	}

	var claims stdClaims
	if err := jsonUnmarshalBase64(claimsBytes, &claims); err != nil {
		return ErrBadJWT
	}

	if claims.Aud != aud {
		return ErrBadJWT
	}

	if claims.Exp < now.Unix() {
		return ErrBadJWT
	}

	if claims.Nbf > now.Unix() {
		return ErrBadJWT
	}

	if err := jsonUnmarshalBase64(claimsBytes, v); err != nil {
		return ErrBadJWT
	}
	return nil
}

func Sign(kid string, priv *ecdsa.PrivateKey, claims any) string {
	header := headerData{
		KID: kid,
		Alg: "ES256",
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		panic(fmt.Errorf("marshal header: %w", err))
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		panic(fmt.Errorf("marshal claims: %w", err))
	}

	var out []byte
	out = base64.RawURLEncoding.AppendEncode(out, headerBytes)
	out = append(out, '.')
	out = base64.RawURLEncoding.AppendEncode(out, claimsBytes)

	hash := sha256.Sum256(out)
	sigR, sigS, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		panic(fmt.Errorf("sign: %w", err))
	}

	// we need to take care to left-0-pad r and s; sig is a buffer to make that easier
	sig := make([]byte, 64)
	r := sigR.Bytes()
	s := sigS.Bytes()
	copy(sig[32-len(r):], r)
	copy(sig[64-len(s):], s)

	out = append(out, '.')
	out = base64.RawURLEncoding.AppendEncode(out, sig)
	return string(out)
}

func parse(token string) ([]byte, []byte, []byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, nil, ErrBadJWT
	}
	return []byte(parts[0]), []byte(parts[1]), []byte(parts[2]), nil
}

func jsonUnmarshalBase64(data []byte, v any) error {
	b, err := base64.RawURLEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
