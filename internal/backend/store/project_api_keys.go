package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) ListProjectAPIKeys(ctx context.Context, req *backendv1.ListProjectAPIKeysRequest) (*backendv1.ListProjectAPIKeysResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	qAPIKeys, err := q.ListProjectAPIKeys(ctx, queries.ListProjectAPIKeysParams{
		ProjectID: projectid.ProjectID(ctx),
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list project API keys: %w", err)
	}

	var projectAPIKeys []*backendv1.ProjectAPIKey
	for _, qAPIKey := range qAPIKeys {
		projectAPIKeys = append(projectAPIKeys, parseProjectAPIKey(qAPIKey))
	}

	var nextPageToken string
	if len(projectAPIKeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(projectAPIKeys[limit].Id)
		projectAPIKeys = projectAPIKeys[:limit]
	}

	return &backendv1.ListProjectAPIKeysResponse{
		ProjectApiKeys: projectAPIKeys,
		NextPageToken:  nextPageToken,
	}, nil
}

func (s *Store) GetProjectAPIKey(ctx context.Context, req *backendv1.GetProjectAPIKeyRequest) (*backendv1.GetProjectAPIKeyResponse, error) {
	id, err := idformat.ProjectAPIKey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse project api key id: %w", err)
	}

	qProjectAPIKey, err := s.q.GetProjectAPIKey(ctx, queries.GetProjectAPIKeyParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        id,
	})
	if err != nil {
		return nil, fmt.Errorf("get project api key: %w", err)
	}

	return &backendv1.GetProjectAPIKeyResponse{ProjectApiKey: parseProjectAPIKey(qProjectAPIKey)}, nil
}

func (s *Store) CreateProjectAPIKey(ctx context.Context, req *backendv1.CreateProjectAPIKeyRequest) (*backendv1.CreateProjectAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	token := uuid.New()
	tokenSHA256 := sha256.Sum256(token[:])
	qProjectAPIKey, err := q.CreateProjectAPIKey(ctx, queries.CreateProjectAPIKeyParams{
		ID:                uuid.New(),
		ProjectID:         projectid.ProjectID(ctx),
		DisplayName:       req.ProjectApiKey.DisplayName,
		SecretTokenSha256: tokenSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create project api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	projectAPIKey := parseProjectAPIKey(qProjectAPIKey)
	projectAPIKey.SecretToken = idformat.ProjectAPIKeySecretToken.Format(token)
	return &backendv1.CreateProjectAPIKeyResponse{ProjectApiKey: projectAPIKey}, nil
}

func (s *Store) UpdateProjectAPIKey(ctx context.Context, req *backendv1.UpdateProjectAPIKeyRequest) (*backendv1.UpdateProjectAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectAPIKeyID, err := idformat.ProjectAPIKey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse project api key id: %w", err)
	}

	qProjectAPIKey, err := q.GetProjectAPIKey(ctx, queries.GetProjectAPIKeyParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        projectAPIKeyID,
	})
	if err != nil {
		return nil, fmt.Errorf("get project api key: %w", err)
	}

	updates := queries.UpdateProjectAPIKeyParams{
		ID:          projectAPIKeyID,
		DisplayName: qProjectAPIKey.DisplayName,
	}

	if req.ProjectApiKey.DisplayName != "" {
		updates.DisplayName = req.ProjectApiKey.DisplayName
	}

	qUpdatedProjectAPIKey, err := q.UpdateProjectAPIKey(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update project api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateProjectAPIKeyResponse{ProjectApiKey: parseProjectAPIKey(qUpdatedProjectAPIKey)}, nil
}

func (s *Store) DeleteProjectAPIKey(ctx context.Context, req *backendv1.DeleteProjectAPIKeyRequest) (*backendv1.DeleteProjectAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectAPIKeyID, err := idformat.ProjectAPIKey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse project api key id: %w", err)
	}

	qProjectAPIKey, err := q.GetProjectAPIKey(ctx, queries.GetProjectAPIKeyParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        projectAPIKeyID,
	})
	if err != nil {
		return nil, fmt.Errorf("get project api key: %w", err)
	}

	if qProjectAPIKey.SecretTokenSha256 != nil {
		return nil, fmt.Errorf("project api key must be revoked before deletion")
	}

	if err := q.DeleteProjectAPIKey(ctx, projectAPIKeyID); err != nil {
		return nil, fmt.Errorf("delete project api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteProjectAPIKeyResponse{}, nil
}

func (s *Store) RevokeProjectAPIKey(ctx context.Context, req *backendv1.RevokeProjectAPIKeyRequest) (*backendv1.RevokeProjectAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectAPIKeyID, err := idformat.ProjectAPIKey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse project api key id: %w", err)
	}

	qProjectAPIKey, err := q.RevokeProjectAPIKey(ctx, projectAPIKeyID)
	if err != nil {
		return nil, fmt.Errorf("revoke project api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.RevokeProjectAPIKeyResponse{ProjectApiKey: parseProjectAPIKey(qProjectAPIKey)}, nil
}

func parseProjectAPIKey(qProjectAPIKey queries.ProjectApiKey) *backendv1.ProjectAPIKey {
	return &backendv1.ProjectAPIKey{
		Id:          idformat.ProjectAPIKey.Format(qProjectAPIKey.ID),
		ProjectId:   idformat.Project.Format(qProjectAPIKey.ProjectID),
		DisplayName: qProjectAPIKey.DisplayName,
		SecretToken: "", // intentionally left blank
		Revoked:     qProjectAPIKey.SecretTokenSha256 == nil,
	}
}
