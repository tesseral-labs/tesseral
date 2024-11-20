package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/store"
)

var ErrInvalidJWTTOken = errors.New("invalid JWT token")
var ErrInvalidSigningMethod = errors.New("invalid signing method")

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

func NewJWT(params NewJWTParams) *JWT {
	return &JWT{
		store: params.Store,
	}
}

func (j *JWT) ParseIntermediateSessionJWT(ctx context.Context, tokenString string) (*IntermediateSessionJWTClaims, error) {
	// TODO: Make this use the project's intermediate session signing key
	signingKey, err := j.store.GetSessionSigningKeyByID(ctx, uuid.New().String())
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &IntermediateSessionJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
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
		return nil, ErrInvalidJWTTOken
	}
	
	return claims, nil
}

func (j *JWT) ParseSessionJWT(ctx context.Context, tokenString string) (*SessionJWTClaims, error) {
	// TODO: Make this use the project's intermediate session signing key
	signingKey, err := j.store.GetSessionSigningKeyByID(ctx, uuid.New().String())
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &SessionJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
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
		return nil, ErrInvalidJWTTOken
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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