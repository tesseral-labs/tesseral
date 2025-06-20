package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListBackendAPIKeys(ctx context.Context, req *backendv1.ListBackendAPIKeysRequest) (*backendv1.ListBackendAPIKeysResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

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
	qAPIKeys, err := q.ListBackendAPIKeys(ctx, queries.ListBackendAPIKeysParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        startID,
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list backend api keys: %w", err)
	}

	var backendAPIKeys []*backendv1.BackendAPIKey
	for _, qAPIKey := range qAPIKeys {
		backendAPIKeys = append(backendAPIKeys, parseBackendAPIKey(qAPIKey))
	}

	var nextPageToken string
	if len(backendAPIKeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qAPIKeys[limit].ID)
		backendAPIKeys = backendAPIKeys[:limit]
	}

	return &backendv1.ListBackendAPIKeysResponse{
		BackendApiKeys: backendAPIKeys,
		NextPageToken:  nextPageToken,
	}, nil
}

func (s *Store) GetBackendAPIKey(ctx context.Context, req *backendv1.GetBackendAPIKeyRequest) (*backendv1.GetBackendAPIKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	id, err := idformat.BackendAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid backend api key id", fmt.Errorf("parse backend api key id: %w", err))
	}

	qBackendAPIKey, err := s.q.GetBackendAPIKey(ctx, queries.GetBackendAPIKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("backend api key not found", fmt.Errorf("get backend api key: %w", err))
		}

		return nil, fmt.Errorf("get backend api key: %w", err)
	}

	return &backendv1.GetBackendAPIKeyResponse{BackendApiKey: parseBackendAPIKey(qBackendAPIKey)}, nil
}

func (s *Store) CreateBackendAPIKey(ctx context.Context, req *backendv1.CreateBackendAPIKeyRequest) (*backendv1.CreateBackendAPIKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if !qProject.EntitledBackendApiKeys {
		return nil, fmt.Errorf("not entitled to backend api keys")
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	token := uuid.New()
	tokenSHA256 := sha256.Sum256(token[:])
	qBackendAPIKey, err := q.CreateBackendAPIKey(ctx, queries.CreateBackendAPIKeyParams{
		ID:                uuid.New(),
		ProjectID:         authn.ProjectID(ctx),
		DisplayName:       req.BackendApiKey.DisplayName,
		SecretTokenSha256: tokenSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create backend api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	backendAPIKey := parseBackendAPIKey(qBackendAPIKey)
	backendAPIKey.SecretToken = idformat.BackendAPIKeySecretToken.Format(token)
	return &backendv1.CreateBackendAPIKeyResponse{BackendApiKey: backendAPIKey}, nil
}

func (s *Store) UpdateBackendAPIKey(ctx context.Context, req *backendv1.UpdateBackendAPIKeyRequest) (*backendv1.UpdateBackendAPIKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	backendAPIKeyID, err := idformat.BackendAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid backend api key id", fmt.Errorf("parse backend api key id: %w", err))
	}

	qBackendAPIKey, err := q.GetBackendAPIKey(ctx, queries.GetBackendAPIKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        backendAPIKeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("backend api key not found", fmt.Errorf("get backend api key: %w", err))
		}

		return nil, fmt.Errorf("get backend api key: %w", err)
	}

	updates := queries.UpdateBackendAPIKeyParams{
		ID:          backendAPIKeyID,
		DisplayName: qBackendAPIKey.DisplayName,
	}

	if req.BackendApiKey.DisplayName != "" {
		updates.DisplayName = req.BackendApiKey.DisplayName
	}

	qUpdatedBackendAPIKey, err := q.UpdateBackendAPIKey(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update backend api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateBackendAPIKeyResponse{BackendApiKey: parseBackendAPIKey(qUpdatedBackendAPIKey)}, nil
}

func (s *Store) DeleteBackendAPIKey(ctx context.Context, req *backendv1.DeleteBackendAPIKeyRequest) (*backendv1.DeleteBackendAPIKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	backendAPIKeyID, err := idformat.BackendAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid backend api key id", fmt.Errorf("parse backend api key id: %w", err))
	}

	qBackendAPIKey, err := q.GetBackendAPIKey(ctx, queries.GetBackendAPIKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        backendAPIKeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("backend api key not found", fmt.Errorf("get backend api key: %w", err))
		}

		return nil, fmt.Errorf("get backend api key: %w", err)
	}

	if qBackendAPIKey.SecretTokenSha256 != nil {
		return nil, apierror.NewFailedPreconditionError("backend api key must be revoked before deletion", fmt.Errorf("backend api key must be revoked before deletion"))
	}

	if err := q.DeleteBackendAPIKey(ctx, backendAPIKeyID); err != nil {
		return nil, fmt.Errorf("delete backend api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteBackendAPIKeyResponse{}, nil
}

func (s *Store) RevokeBackendAPIKey(ctx context.Context, req *backendv1.RevokeBackendAPIKeyRequest) (*backendv1.RevokeBackendAPIKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	backendAPIKeyID, err := idformat.BackendAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid backend api key id", fmt.Errorf("parse backend api key id: %w", err))
	}

	qBackendAPIKey, err := q.RevokeBackendAPIKey(ctx, backendAPIKeyID)
	if err != nil {
		return nil, fmt.Errorf("revoke backend api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.RevokeBackendAPIKeyResponse{BackendApiKey: parseBackendAPIKey(qBackendAPIKey)}, nil
}

func parseBackendAPIKey(qBackendAPIKey queries.BackendApiKey) *backendv1.BackendAPIKey {
	return &backendv1.BackendAPIKey{
		Id:          idformat.BackendAPIKey.Format(qBackendAPIKey.ID),
		DisplayName: qBackendAPIKey.DisplayName,
		CreateTime:  timestamppb.New(*qBackendAPIKey.CreateTime),
		UpdateTime:  timestamppb.New(*qBackendAPIKey.UpdateTime),
		SecretToken: "", // intentionally left blank
		Revoked:     qBackendAPIKey.SecretTokenSha256 == nil,
	}
}
