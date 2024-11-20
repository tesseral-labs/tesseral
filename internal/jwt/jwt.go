package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/store"
)

var ErrInvalidJWTFormat = errors.New("invalid JWT format")
var ErrInvalidJWTToken = errors.New("invalid JWT token")
var ErrInvalidSigningMethod = errors.New("invalid signing method")
var ErrKidNotFound = errors.New("`kid` not found in JWT header")

type JWT struct {
	store *store.Store
}

type NewJWTParams struct {
	Store *store.Store
}

type IntermediateSessionJWTClaims struct {
	jwt.Claims

	Email string
	ExpiresAt int64
	IssuedAt int64
	ProjectID string
	Subject string
}

type SessionJWTClaims struct {
	jwt.Claims

	Email string
	ExpiresAt int64
	IssuedAt int64
	OrganizationID string
	ProjectID string
	Subject string
	UserID string
}

func New(params NewJWTParams) *JWT {
	return &JWT{
		store: params.Store,
	}
}

func (j *JWT) ParseIntermediateSessionJWT(ctx context.Context, tokenString string) (*IntermediateSessionJWTClaims, error) {
	kid, err := extractKidFromJWT(tokenString)
	if err != nil {
		return nil, err
	}

	signingKey, err := j.store.GetSessionSigningKeyByID(ctx, kid)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &IntermediateSessionJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return signingKey.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract the claims
	claims, ok := token.Claims.(*IntermediateSessionJWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidJWTToken
	}
	
	return claims, nil
}

func (j *JWT) ParseSessionJWT(ctx context.Context, tokenString string) (*SessionJWTClaims, error) {
	kid, err := extractKidFromJWT(tokenString)
	if err != nil {
		return nil, err
	}

	signingKey, err := j.store.GetSessionSigningKeyByID(ctx, kid)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &SessionJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract the claims
	claims, ok := token.Claims.(*SessionJWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidJWTToken
	}

	return claims, nil
}

func (j *JWT) SignIntermediateSessionJWT(ctx context.Context, claims *IntermediateSessionJWTClaims) (string, error) {
	// TODO: Make this use the project's intermediate session signing key
	signingKey, err := j.store.GetIntermediateSessionSigningKeyByID(ctx, uuid.New().String())
	if err != nil {
		return "", err
	}

	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	token.Header["kid"] = signingKey.ID

	tokenString, err := token.SignedString(signingKey.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWT) SignSessionJWT(ctx context.Context, claims *SessionJWTClaims) (string, error) {
	// TODO: Make this use the project's intermediate session signing key
	signingKey, err := j.store.GetIntermediateSessionSigningKeyByID(ctx, uuid.New().String())
	if err != nil {
		return "", err
	}

	claims.IssuedAt = time.Now().Unix()
	// TODO: Make this honor the project's activity timeout
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.Header["kid"] = signingKey.ID

	tokenString, err :=  token.SignedString(signingKey.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func extractKidFromJWT(tokenString string) (string, error) {
	// Split the token into its three parts: Header, Payload, and Signature
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", ErrInvalidJWTFormat
	}

	// Decode the header (first part)
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	// Parse the header as JSON
	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return "", err
	}

	// Extract the `kid` field
	kid, ok := header["kid"].(string)
	if !ok {
		return "", ErrKidNotFound
	}

	return kid, nil
}