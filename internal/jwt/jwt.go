package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidJWTTOken = errors.New("invalid JWT token")

type JWT struct {}

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

func ParseIntermediateSessionJWT(tokenString string) (*IntermediateSessionJWTClaims, error) {
	// TODO: Make this use the project's intermediate session signing key
	signingKey := []byte("")

	token, err := jwt.ParseWithClaims(tokenString, &IntermediateSessionJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
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

func ParseSessionJWT(tokenString string) (*SessionJWTClaims, error) {
	// TODO: Make this use the project's session signing key
	signingKey := []byte("")

	token, err := jwt.ParseWithClaims(tokenString, &SessionJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
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

func SignIntermediateSessionJWT(ctx context.Context, claims *IntermediateSessionJWTClaims) (string, error) {
	// TODO: Make this use the project's intermediate session signing key
	signingKey := []byte("")

	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func SignSessionJWT(ctx context.Context, claims *SessionJWTClaims) (string, error) {
	// TODO: Make this use the project's session signing key
	signingKey := []byte("")

	claims.IssuedAt = time.Now().Unix()
	// TODO: Make this honor the project's activity timeout
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}