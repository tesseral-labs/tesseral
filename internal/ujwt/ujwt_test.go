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

	"github.com/tesseral-labs/tesseral/internal/ujwt"
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

func TestClaims_invalid(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "bad number of parts",
			token: ".",
		},
		{
			name:  "invalid base64 header",
			token: "invalidbase64..",
		},
		{
			// echo '{"alg":"ES256"}' | base64
			name:  "invalid base64 claims",
			token: "eyJhbGciOiJFUzI1NiJ9Cg.invalidbase64.",
		},
		{
			// echo '{"alg":"ES256"}' | base64
			name:  "invalid base64 signature",
			token: "eyJhbGciOiJFUzI1NiJ9Cg..invalidbase64",
		},
		{
			// echo '{"alg":"HS256"}' | base64
			name:  "bad alg",
			token: "eyJhbGciOiJIUzI1NiJ9Cg..",
		},
		{
			// echo '{' | base64
			name:  "invalid JSON header",
			token: "ewo..",
		},
		{
			// echo '{' | base64
			name:  "invalid JSON claims",
			token: ".ewo.",
		},
		{
			// echo '{"alg":"ES256"}' | base64
			// echo '{}' | base64
			// echo -n 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' | base64
			name:  "incorrect ecdsa signature correct length",
			token: "eyJhbGciOiJFUzI1NiJ9Cg.e30K.YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYQ",
		},
		{
			// aud=aud2 nbf=1 exp=3
			name:  "bad aud",
			token: "eyJraWQiOiJhYWEiLCJhbGciOiJFUzI1NiJ9.eyJhdWQiOiJhdWQyIiwiZXhwIjozLCJuYmYiOjF9.xCSSUjgFCnXbO3XnHIJCAQTegr4CWKxGgSkXBIqZMhFupnXf6thm4itkCZwX_7QbM28y25f0m09Zyg4llVEVNg",
		},
		{
			// aud=aud1 nbf=3 exp=3
			name:  "bad nbf",
			token: "eyJraWQiOiJhYWEiLCJhbGciOiJFUzI1NiJ9.eyJhdWQiOiJhdWQxIiwiZXhwIjozLCJuYmYiOjN9.3IpKkdPhUkBnOvjkSRo6OIB6Ijli9iddiFTV6u-lCqU1Jn5PV8yLtUxauWDfY3ejnNxGgH_Pet5iPo6FQL4ULg",
		},
		{
			// aud=aud1 nbf=1 exp=1
			name:  "bad exp",
			token: "eyJraWQiOiJhYWEiLCJhbGciOiJFUzI1NiJ9.eyJhdWQiOiJhdWQxIiwiZXhwIjoxLCJuYmYiOjF9.W8SR1l9UdCi3qbWY2OvY7vgm0e-qEhgu24vL_llAMM2NXM6eSalJc6mhX24leoft29jcVk1Q-YQOnFAz7m7A5g",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var claims map[string]any
			err := ujwt.Claims(&priv.PublicKey, "aud1", time.Unix(2, 0), &claims, tt.token)
			if !errors.Is(err, ujwt.ErrBadJWT) {
				t.Errorf("ujwt.Claims() error = %v, wantErr %v", err, ujwt.ErrBadJWT)
			}
		})
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
