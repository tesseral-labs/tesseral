package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) GetSessionSigningKeyPublicKey(ctx context.Context, sessionSigningKeyID string) (*ecdsa.PublicKey, error) {
	sessionSigningKeyUUID, err := idformat.SessionSigningKey.Parse(sessionSigningKeyID)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := s.q.GetSessionSigningKeyPublicKey(ctx, queries.GetSessionSigningKeyPublicKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        sessionSigningKeyUUID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("session signing key not found", fmt.Errorf("get session signing key public key: %w", err))
		}

		return nil, fmt.Errorf("get session signing key public key: %w", err)
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	return publicKey.(*ecdsa.PublicKey), nil
}
