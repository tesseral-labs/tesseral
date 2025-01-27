package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
)

type GetSessionPublicKeysByProjectIDResponseKey struct {
	ID        string
	PublicKey *ecdsa.PublicKey
}

func (s *Store) GetSessionPublicKeysByProjectID(ctx context.Context, projectId string) ([]GetSessionPublicKeysByProjectIDResponseKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(projectId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid project id", fmt.Errorf("parse project id: %w", err))
	}

	sessionSigningKeys, err := q.GetSessionSigningKeysByProjectID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("session signing keys not found", fmt.Errorf("get session signing keys by project id: %w", err))
		}

		return nil, fmt.Errorf("get session signing keys by project id: %w", err)
	}

	var out []GetSessionPublicKeysByProjectIDResponseKey
	for _, sessionSigningKey := range sessionSigningKeys {
		pub, err := x509.ParsePKIXPublicKey(sessionSigningKey.PublicKey)
		if err != nil {
			panic(fmt.Errorf("public key from bytes: %w", err))
		}

		out = append(out, GetSessionPublicKeysByProjectIDResponseKey{
			ID:        idformat.SessionSigningKey.Format(sessionSigningKey.ID),
			PublicKey: pub.(*ecdsa.PublicKey),
		})
	}

	return out, nil
}
