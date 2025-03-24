package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

var ErrBadBackendAPIKey = fmt.Errorf("bad backend api key")

type AuthenticateBackendAPIKeyResponse struct {
	BackendAPIKeyID string
	ProjectID       string
}

func (s *Store) AuthenticateBackendAPIKey(ctx context.Context, bearerToken string) (*AuthenticateBackendAPIKeyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	secretToken, err := idformat.BackendAPIKeySecretToken.Parse(bearerToken)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid backend api key secret token", fmt.Errorf("parse backend api key secret token: %w", err))
	}

	secretTokenSHA := sha256.Sum256(secretToken[:])
	qBackendAPIKey, err := q.GetBackendAPIKeyBySecretTokenSHA256(ctx, secretTokenSHA[:])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("backend api key not found", fmt.Errorf("backend api key not found"))
		}

		return nil, fmt.Errorf("get backend api key by secret token sha256: %w", err)
	}

	return &AuthenticateBackendAPIKeyResponse{
		BackendAPIKeyID: idformat.BackendAPIKey.Format(qBackendAPIKey.ID),
		ProjectID:       idformat.Project.Format(qBackendAPIKey.ProjectID),
	}, nil
}
