package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/wellknown/authn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateAPIKey(ctx context.Context, req *backendv1.CreateAPIKeyRequest) (*backendv1.CreateAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	if _, err := s.q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ID:        orgID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization: %w", err))
		}
		return nil, fmt.Errorf("get organization: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if !qProject.ApiKeysEnabled {
		return nil, apierror.NewPermissionDeniedError("api keys are not enabled for this project", fmt.Errorf("api keys not enabled for project"))
	}

	secretTokenUUID := uuid.New()
	secretTokenFormatter := idformat.APIKey

	// Handle custom api key prefixes
	if qProject.ApiKeysPrefix != nil {
		secretTokenFormatter = idformat.MustNewFormat(*qProject.ApiKeysPrefix)
	}

	secretToken := secretTokenFormatter.Format(secretTokenUUID)

	qAPIKey, err := q.CreateAPIKey(ctx, queries.CreateAPIKeyParams{
		ID:             uuid.New(),
		DisplayName:    req.DisplayName,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.CreateAPIKeyResponse{
		ApiKey: parseAPIKey(qAPIKey, &secretToken),
	}, nil
}

func (s *Store) DeleteAPIKey(ctx context.Context, req *backendv1.DeleteAPIKeyRequest) (*backendv1.DeleteAPIKeyResponse, error) {
	return &backendv1.DeleteAPIKeyResponse{}, nil
}

func (s *Store) GetAPIKey(ctx context.Context, req *backendv1.GetAPIKeyRequest) (*backendv1.GetAPIKeyResponse, error) {
	return &backendv1.GetAPIKeyResponse{}, nil
}

func (s *Store) ListAPIKeys(ctx context.Context, req *backendv1.ListAPIKeysRequest) (*backendv1.ListAPIKeysResponse, error) {
	return &backendv1.ListAPIKeysResponse{}, nil
}

func (s *Store) RevokeAPIKey(ctx context.Context, req *backendv1.RevokeAPIKeyRequest) (*backendv1.RevokeAPIKeyResponse, error) {
	return &backendv1.RevokeAPIKeyResponse{}, nil
}

func (s *Store) UpdateAPIKey(ctx context.Context, req *backendv1.UpdateAPIKeyRequest) (*backendv1.UpdateAPIKeyResponse, error) {
	return &backendv1.UpdateAPIKeyResponse{}, nil
}

func parseAPIKey(qAPIKey queries.ApiKey, secretToken *string) *backendv1.APIKey {
	return &backendv1.APIKey{
		Id:                idformat.APIKey.Format(qAPIKey.ID),
		DisplayName:       qAPIKey.DisplayName,
		CreateTime:        timestamppb.New(*qAPIKey.CreateTime),
		UpdateTime:        timestamppb.New(*qAPIKey.UpdateTime),
		ExpireTime:        timestamppb.New(*qAPIKey.ExpireTime),
		Revoked:           qAPIKey.SecretTokenSha256 == nil,
		SecretToken:       derefOrEmpty(secretToken),
		SecretTokenSuffix: derefOrEmpty(qAPIKey.SecretTokenSuffix),
	}
}
