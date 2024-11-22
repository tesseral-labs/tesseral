package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

func (s *Store) CreateProjectAPIKey(ctx context.Context, req *backendv1.CreateProjectAPIKeyRequest) (*backendv1.CreateProjectAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// todo project ID should be determined from authn ctx
	projectID, err := idformat.Project.Parse(req.ProjectApiKey.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("parse project id: %w", err)
	}

	secretToken := uuid.New()
	secretTokenSHA := sha256.Sum256(secretToken[:])
	qProjectAPIKey, err := q.CreateProjectAPIKey(ctx, queries.CreateProjectAPIKeyParams{
		ID:                uuid.New(),
		ProjectID:         projectID,
		SecretTokenSha256: secretTokenSHA[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create project api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	projectAPIKey := parseProjectAPIKey(qProjectAPIKey)
	projectAPIKey.SecretToken = idformat.ProjectAPIKeySecretToken.Format(secretToken)
	return &backendv1.CreateProjectAPIKeyResponse{ProjectApiKey: projectAPIKey}, nil
}

func parseProjectAPIKey(qProjectAPIKey queries.ProjectApiKey) *backendv1.ProjectAPIKey {
	return &backendv1.ProjectAPIKey{
		Id:          idformat.ProjectAPIKey.Format(qProjectAPIKey.ID),
		ProjectId:   idformat.Project.Format(qProjectAPIKey.ProjectID),
		SecretToken: "", // intentionally left blank
	}
}
