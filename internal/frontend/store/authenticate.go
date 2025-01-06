package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"time"

	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetSessionSigningKeyPublicKey(ctx context.Context, sessionSigningKeyID string) (*ecdsa.PublicKey, error) {
	sessionSigningKeyUUID, err := idformat.SessionSigningKey.Parse(sessionSigningKeyID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	publicKeyBytes, err := s.q.GetSessionSigningKeyPublicKey(ctx, queries.GetSessionSigningKeyPublicKeyParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        sessionSigningKeyUUID,
		Now:       &now,
	})
	if err != nil {
		return nil, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*ecdsa.PublicKey), nil
}
