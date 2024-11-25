package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
)

var ErrBadProjectAPIKey = fmt.Errorf("bad project api key")

func (s *Store) AuthenticateProjectAPIKey(ctx context.Context, bearerToken string) (*backendv1.ProjectAPIKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	secretToken, err := idformat.ProjectAPIKeySecretToken.Parse(bearerToken)
	if err != nil {
		return nil, fmt.Errorf("parse project api key secret token: %w", err)
	}

	secretTokenSHA := sha256.Sum256(secretToken[:])
	qProjectAPIKey, err := q.GetProjectAPIKeyBySecretTokenSHA256(ctx, secretTokenSHA[:])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBadProjectAPIKey
		}

		return nil, fmt.Errorf("get project api key by secret token sha256: %w", err)
	}

	return parseProjectAPIKey(qProjectAPIKey), nil
}
