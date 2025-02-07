package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
)

var ErrBadProjectAPIKey = fmt.Errorf("bad project api key")

func (s *Store) AuthenticateProjectAPIKey(ctx context.Context, bearerToken string) (*backendv1.ProjectAPIKey, *backendv1.Project, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer rollback()

	secretToken, err := idformat.ProjectAPIKeySecretToken.Parse(bearerToken)
	if err != nil {
		return nil, nil, apierror.NewInvalidArgumentError("invalid project api key secret token", fmt.Errorf("parse project api key secret token: %w", err))
	}

	secretTokenSHA := sha256.Sum256(secretToken[:])
	qProjectAPIKey, err := q.GetProjectAPIKeyBySecretTokenSHA256(ctx, secretTokenSHA[:])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, apierror.NewNotFoundError("project api key not found", fmt.Errorf("project api key not found"))
		}

		return nil, nil, fmt.Errorf("get project api key by secret token sha256: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, qProjectAPIKey.ProjectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, nil, fmt.Errorf("get project by id: %w", err)
	}

	qProjectPasskeyRPIDs, err := q.GetProjectPasskeyRPIDs(ctx, qProject.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("get project passkey rp ids: %w", err)
	}

	return parseProjectAPIKey(qProjectAPIKey), parseProject(&qProject, qProjectPasskeyRPIDs), nil
}
