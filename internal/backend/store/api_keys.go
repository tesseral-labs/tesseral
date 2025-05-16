package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/store/secretformat"
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

	var secretTokenValue [35]byte
	if _, err := rand.Read(secretTokenValue[:]); err != nil {
		return nil, fmt.Errorf("generate secret token: %w", err)
	}
	secretTokenFormatter := secretformat.APIKeySecretToken

	// Handle custom api key prefixes
	if qProject.ApiKeysPrefix != nil {
		secretTokenFormatter = secretformat.MustNewFormat(*qProject.ApiKeysPrefix)
	}

	secretToken := secretTokenFormatter.Format(secretTokenValue)
	secretTokenSuffix := secretToken[len(secretToken)-5:]

	secretTokenSha256 := sha256.Sum256(secretTokenValue[:])

	var expireTime *time.Time
	if req.ExpireTime != nil {
		formattedExpireTime := req.ExpireTime.AsTime()
		expireTime = &formattedExpireTime
	}

	qAPIKey, err := q.CreateAPIKey(ctx, queries.CreateAPIKeyParams{
		ID:                uuid.New(),
		DisplayName:       req.DisplayName,
		ExpireTime:        expireTime,
		OrganizationID:    orgID,
		SecretTokenSha256: secretTokenSha256[:],
		SecretTokenSuffix: &secretTokenSuffix,
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
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyID, err := idformat.APIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	if err := q.DeleteAPIKey(ctx, queries.DeleteAPIKeyParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("delete api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.DeleteAPIKeyResponse{}, nil
}

func (s *Store) GetAPIKey(ctx context.Context, req *backendv1.GetAPIKeyRequest) (*backendv1.GetAPIKeyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}

	defer rollback()

	apiKeyID, err := idformat.APIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	qAPIKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get api key: %w", err))
		}

		return nil, fmt.Errorf("get api key: %w", err)
	}

	return &backendv1.GetAPIKeyResponse{
		ApiKey: parseAPIKey(qAPIKey, nil),
	}, nil
}

func (s *Store) ListAPIKeys(ctx context.Context, req *backendv1.ListAPIKeysRequest) (*backendv1.ListAPIKeysResponse, error) {
	return &backendv1.ListAPIKeysResponse{}, nil
}

func (s *Store) RevokeAPIKey(ctx context.Context, req *backendv1.RevokeAPIKeyRequest) (*backendv1.RevokeAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyID, err := idformat.APIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	if err := q.RevokeAPIKey(ctx, queries.RevokeAPIKeyParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("revoke api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.RevokeAPIKeyResponse{}, nil
}

func (s *Store) UpdateAPIKey(ctx context.Context, req *backendv1.UpdateAPIKeyRequest) (*backendv1.UpdateAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyID, err := idformat.APIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	updatedApiKey, err := q.UpdateAPIKey(ctx, queries.UpdateAPIKeyParams{
		ID:          apiKeyID,
		DisplayName: req.DisplayName,
		ProjectID:   authn.ProjectID(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("update api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.UpdateAPIKeyResponse{
		ApiKey: parseAPIKey(updatedApiKey, nil),
	}, nil
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
