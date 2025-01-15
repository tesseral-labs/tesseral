package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/errorcodes"
	"github.com/openauth/openauth/internal/frontend/authn"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetSessionSigningKeyPublicKey(ctx context.Context, sessionSigningKeyID string) (*ecdsa.PublicKey, error) {
	sessionSigningKeyUUID, err := idformat.SessionSigningKey.Parse(sessionSigningKeyID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	publicKeyBytes, err := s.q.GetSessionSigningKeyPublicKey(ctx, queries.GetSessionSigningKeyPublicKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        sessionSigningKeyUUID,
		Now:       &now,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errorcodes.NewNotFoundError())
		}

		return nil, fmt.Errorf("get session signing key public key: %w", err)
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	return publicKey.(*ecdsa.PublicKey), nil
}
