package ujwt_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/openauth-dev/openauth/internal/ujwt"
)

func TestKeyID(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    string
		wantErr error
	}{
		{
			// echo '{"alg":"ES256","kid":"foobar"}' | base64
			name:    "valid header",
			token:   "eyJhbGciOiJFUzI1NiIsImtpZCI6ImZvb2JhciJ9Cg..",
			want:    "foobar",
			wantErr: nil,
		},
		{
			name:    "bad number of parts",
			token:   ".",
			want:    "",
			wantErr: ujwt.ErrBadJWT,
		},
		{
			name:    "invalid base64 header",
			token:   "invalidbase64..",
			want:    "",
			wantErr: ujwt.ErrBadJWT,
		},
		{
			// echo '{"alg":"HS256"}' | base64
			name:    "bad alg",
			token:   "eyJhbGciOiJIUzI1NiJ9Cg..",
			want:    "",
			wantErr: ujwt.ErrBadJWT,
		},
		{
			// echo '{' | base64
			name:    "invalid JSON header",
			token:   "ewo..",
			want:    "",
			wantErr: ujwt.ErrBadJWT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ujwt.KeyID(tt.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("KeyID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KeyID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// privKeyPEM stores an insecure PEM-encoded ECDSA private key for predictable
// unit tests.
//
// generated with: openssl ecparam -genkey -name prime256v1 -noout -out mykey.pem
const (
	privKeyPEM = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKQVPEgRbSYPm7OYcdGGe4dSnip0hhAV4l9U7KfRADZYoAoGCCqGSM49
AwEHoUQDQgAEukQA0ieZ9LT61RMSBgCq5nrrYOJwLjdkQIHroU6inSI9Xy9UYNJk
12mcoT8grb6F9Kvac9KhPjzjpeSl72CapA==
-----END EC PRIVATE KEY-----
`
)

var (
	priv *ecdsa.PrivateKey
)

func init() {
	block, _ := pem.Decode([]byte(privKeyPEM))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		panic(fmt.Errorf("bad ec pem block"))
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("parse ec private key: %w", err))
	}

	priv = privateKey
}

func TestSign(t *testing.T) {
	r := rand.Reader
	rand.Reader = zeroReader{}
	defer func() { rand.Reader = r }()

	jwt := ujwt.Sign("aaa", priv, map[string]any{"aud": "aud1", "nbf": 1, "exp": 3})

	wantJWT := "eyJraWQiOiJhYWEiLCJhbGciOiJFUzI1NiJ9.eyJhdWQiOiJhdWQxIiwiZXhwIjozLCJuYmYiOjF9.zGUJGj1SBoNmi9A1MmJd-QSA4ri8Fnr4Y1oOEMl2D6XDvdSXb2SVdmGmO3kdJHwn90dEDyFACdL1F3-wy2G0Gg"
	if jwt != wantJWT {
		t.Errorf("ujwt.Sign() = %v, want %v", jwt, wantJWT)
	}
}

func TestClaims(t *testing.T) {
	jwt := "eyJraWQiOiJhYWEiLCJhbGciOiJFUzI1NiJ9.eyJhdWQiOiJhdWQxIiwiZXhwIjozLCJuYmYiOjF9.zGUJGj1SBoNmi9A1MmJd-QSA4ri8Fnr4Y1oOEMl2D6XDvdSXb2SVdmGmO3kdJHwn90dEDyFACdL1F3-wy2G0Gg"

	var claims map[string]any
	if err := ujwt.Claims(&priv.PublicKey, "aud1", time.Unix(2, 0), &claims, jwt); err != nil {
		t.Fatalf("ujwt.Claims() error = %v", err)
	}

	wantClaims := map[string]any{
		"aud": "aud1",
		"nbf": 1.0,
		"exp": 3.0,
	}

	if !reflect.DeepEqual(claims, wantClaims) {
		t.Errorf("ujwt.Claims() = %v, want %v", claims, wantClaims)
	}
}

// zeroReader is an insecure crypto/rand.Reader for predictable unit tests.
type zeroReader struct{}

func (zeroReader) Read(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf), nil
}
