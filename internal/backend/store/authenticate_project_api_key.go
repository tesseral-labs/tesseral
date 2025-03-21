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

var ErrBadProjectAPIKey = fmt.Errorf("bad project api key")

type AuthenticateProjectAPIKeyResponse struct {
	ProjectAPIKeyID string
	ProjectID       string
}

func (s *Store) AuthenticateProjectAPIKey(ctx context.Context, bearerToken string) (*AuthenticateProjectAPIKeyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	secretToken, err := idformat.ProjectAPIKeySecretToken.Parse(bearerToken)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid project api key secret token", fmt.Errorf("parse project api key secret token: %w", err))
	}

	secretTokenSHA := sha256.Sum256(secretToken[:])
	qProjectAPIKey, err := q.GetProjectAPIKeyBySecretTokenSHA256(ctx, secretTokenSHA[:])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project api key not found", fmt.Errorf("project api key not found"))
		}

		return nil, fmt.Errorf("get project api key by secret token sha256: %w", err)
	}

	return &AuthenticateProjectAPIKeyResponse{
		ProjectAPIKeyID: idformat.ProjectAPIKey.Format(qProjectAPIKey.ID),
		ProjectID:       idformat.Project.Format(qProjectAPIKey.ProjectID),
	}, nil
}
