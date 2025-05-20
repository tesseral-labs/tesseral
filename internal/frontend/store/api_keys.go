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
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/prettysecret"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const apiKeySecretTokenSuffixLength = 4

func (s *Store) CreateAPIKey(ctx context.Context, req *frontendv1.CreateAPIKeyRequest) (*frontendv1.CreateAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if _, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx)); err != nil {
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

	if qProject.ApiKeySecretTokenPrefix == nil || *qProject.ApiKeySecretTokenPrefix == "" {
		return nil, apierror.NewInvalidArgumentError("api key secret token prefix is required", fmt.Errorf("api key secret token prefix is required"))
	}

	var secretTokenValue [35]byte
	if _, err := rand.Read(secretTokenValue[:]); err != nil {
		return nil, fmt.Errorf("generate secret token: %w", err)
	}

	// Handle custom api key prefixes
	secretToken := prettysecret.Format(*qProject.ApiKeySecretTokenPrefix, secretTokenValue)
	secretTokenSuffix := secretToken[len(secretToken)-apiKeySecretTokenSuffixLength:]

	secretTokenSha256 := sha256.Sum256(secretTokenValue[:])

	var expireTime *time.Time
	if req.ApiKey.ExpireTime != nil {
		formattedExpireTime := req.ApiKey.ExpireTime.AsTime()
		expireTime = &formattedExpireTime
	}

	qAPIKey, err := q.CreateAPIKey(ctx, queries.CreateAPIKeyParams{
		ID:                uuid.New(),
		DisplayName:       req.ApiKey.DisplayName,
		ExpireTime:        expireTime,
		OrganizationID:    authn.OrganizationID(ctx),
		SecretTokenSha256: secretTokenSha256[:],
		SecretTokenSuffix: &secretTokenSuffix,
	})
	if err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.CreateAPIKeyResponse{
		ApiKey: parseAPIKey(qAPIKey, &secretToken),
	}, nil
}

func (s *Store) DeleteAPIKey(ctx context.Context, req *frontendv1.DeleteAPIKeyRequest) (*frontendv1.DeleteAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyID, err := idformat.APIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	qApiKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:             apiKeyID,
		OrganizationID: authn.OrganizationID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get api key: %w", err))
		}
		return nil, fmt.Errorf("get api key: %w", err)
	}

	if qApiKey.SecretTokenSha256 != nil {
		return nil, apierror.NewFailedPreconditionError("api key must be revoked to be deleted", fmt.Errorf("api key mut be revoked to be deleted"))
	}

	if err := q.DeleteAPIKey(ctx, queries.DeleteAPIKeyParams{
		ID:             apiKeyID,
		OrganizationID: authn.OrganizationID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("delete api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.DeleteAPIKeyResponse{}, nil
}

func (s *Store) GetAPIKey(ctx context.Context, req *frontendv1.GetAPIKeyRequest) (*frontendv1.GetAPIKeyResponse, error) {
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
		ID:             apiKeyID,
		OrganizationID: authn.OrganizationID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get api key: %w", err))
		}

		return nil, fmt.Errorf("get api key: %w", err)
	}

	return &frontendv1.GetAPIKeyResponse{
		ApiKey: parseAPIKey(qAPIKey, nil),
	}, nil
}

func (s *Store) ListAPIKeys(ctx context.Context, req *frontendv1.ListAPIKeysRequest) (*frontendv1.ListAPIKeysResponse, error) {
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
	qAPIKeys, err := q.ListAPIKeys(ctx, queries.ListAPIKeysParams{
		ID:             startID,
		OrganizationID: authn.OrganizationID(ctx),
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}

	apiKeys := make([]*frontendv1.APIKey, len(qAPIKeys))
	for i, qAPIKey := range qAPIKeys {
		apiKeys[i] = parseAPIKey(qAPIKey, nil)
	}

	var nextPageToken string
	if len(apiKeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qAPIKeys[limit].ID)
		apiKeys = apiKeys[:limit]
	}

	return &frontendv1.ListAPIKeysResponse{
		ApiKeys:       apiKeys,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) RevokeAPIKey(ctx context.Context, req *frontendv1.RevokeAPIKeyRequest) (*frontendv1.RevokeAPIKeyResponse, error) {
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
		ID:             apiKeyID,
		OrganizationID: authn.OrganizationID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("revoke api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.RevokeAPIKeyResponse{}, nil
}

func (s *Store) UpdateAPIKey(ctx context.Context, req *frontendv1.UpdateAPIKeyRequest) (*frontendv1.UpdateAPIKeyResponse, error) {
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
		ID:             apiKeyID,
		DisplayName:    req.ApiKey.DisplayName,
		OrganizationID: authn.OrganizationID(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("update api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.UpdateAPIKeyResponse{
		ApiKey: parseAPIKey(updatedApiKey, nil),
	}, nil
}

func parseAPIKey(qAPIKey queries.ApiKey, secretToken *string) *frontendv1.APIKey {
	return &frontendv1.APIKey{
		Id:                idformat.APIKey.Format(qAPIKey.ID),
		DisplayName:       qAPIKey.DisplayName,
		CreateTime:        timestamppb.New(*qAPIKey.CreateTime),
		UpdateTime:        timestamppb.New(*qAPIKey.UpdateTime),
		ExpireTime:        timestampOrNil(qAPIKey.ExpireTime),
		Revoked:           qAPIKey.SecretTokenSha256 == nil,
		SecretToken:       derefOrEmpty(secretToken),
		SecretTokenSuffix: derefOrEmpty(qAPIKey.SecretTokenSuffix),
	}
}
