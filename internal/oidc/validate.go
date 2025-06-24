package oidc

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type IDTokenClaims struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Email string `json:"email"`
}

// fetchJWKS fetches the JSON Web Key Set from the specified URI.
func (c *Client) fetchJWKS(ctx context.Context, jwksURI string) (*jwks, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for JWKS: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code fetching JWKS: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response body: %w", err)
	}

	var jwks jwks
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWKS: %w", err)
	}

	return &jwks, nil
}

// findKey finds the matching JWK in the JWKS based on the key ID (kid).
func findKey(kid string, jwks *jwks) (*jwk, error) {
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return &key, nil
		}
	}
	return nil, fmt.Errorf("key with kid '%s' not found in JWKS", kid)
}

// verifyRSASignature verifies the token's signature using the provided public key.
func verifyRSASignature(tokenString string, jwk *jwk) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return errors.New("token format is invalid, expected 3 parts")
	}

	signedContent := parts[0] + "." + parts[1]

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return fmt.Errorf("failed to decode modulus (n): %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return fmt.Errorf("failed to decode exponent (e): %w", err)
	}

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}

	hashed := sha256.Sum256([]byte(signedContent))

	// RS256 uses PKCS1v15 padding.
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

type ValidateIDTokenRequest struct {
	IDToken       string
	Issuer        string
	Configuration *Configuration
}

// ValidateIDToken validates an ID token and returns the claims if valid.
func (c *Client) ValidateIDToken(ctx context.Context, req ValidateIDTokenRequest) (*IDTokenClaims, error) {
	parts := strings.Split(req.IDToken, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token header: %w", err)
	}
	var header tokenHeader
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token header: %w", err)
	}

	jwks, err := c.fetchJWKS(ctx, req.Configuration.JwksURI)
	if err != nil {
		return nil, fmt.Errorf("could not fetch JWKS: %w", err)
	}

	key, err := findKey(header.Kid, jwks)
	if err != nil {
		return nil, fmt.Errorf("could not find matching key: %w", err)
	}

	// TODO: Support for other algorithms can be added later
	if key.Alg != "RS256" {
		return nil, fmt.Errorf("unsupported key algorithm: %s, expected RS256", key.Alg)
	}

	if err := verifyRSASignature(req.IDToken, key); err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token claims: %w", err)
	}
	var claims IDTokenClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token claims: %w", err)
	}

	if claims.Iss != req.Issuer {
		return nil, fmt.Errorf("issuer mismatch: expected %s, got %s", req.Issuer, claims.Iss)
	}
	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return &claims, nil
}
