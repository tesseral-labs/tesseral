package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListSCIMAPIKeys(ctx context.Context, req *frontendv1.ListSCIMAPIKeysRequest) (*frontendv1.ListSCIMAPIKeysResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qSCIMAPIKeys, err := q.ListSCIMAPIKeys(ctx, queries.ListSCIMAPIKeysParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list scim api keys: %w", err)
	}

	var scimAPIKeys []*frontendv1.SCIMAPIKey
	for _, qSCIMAPIKey := range qSCIMAPIKeys {
		scimAPIKeys = append(scimAPIKeys, parseSCIMAPIKey(qSCIMAPIKey))
	}

	var nextPageToken string
	if len(scimAPIKeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qSCIMAPIKeys[limit].ID)
		scimAPIKeys = scimAPIKeys[:limit]
	}

	return &frontendv1.ListSCIMAPIKeysResponse{
		ScimApiKeys:   scimAPIKeys,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetSCIMAPIKey(ctx context.Context, req *frontendv1.GetSCIMAPIKeyRequest) (*frontendv1.GetSCIMAPIKeyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	scimAPIKeyID, err := idformat.SCIMAPIKey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse scim api key id: %w", err)
	}

	qSCIMAPIKey, err := q.GetSCIMAPIKey(ctx, queries.GetSCIMAPIKeyParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             scimAPIKeyID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apierror.NewNotFoundError("scim api key not found", fmt.Errorf("get scim api key: %w", err))
		}

		return nil, fmt.Errorf("get scim api key: %w", err)
	}

	return &frontendv1.GetSCIMAPIKeyResponse{ScimApiKey: parseSCIMAPIKey(qSCIMAPIKey)}, nil
}

func (s *Store) CreateSCIMAPIKey(ctx context.Context, req *frontendv1.CreateSCIMAPIKeyRequest) (*frontendv1.CreateSCIMAPIKeyResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// authz
	qOrg, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apierror.NewFailedPreconditionError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	if !qOrg.ScimEnabled {
		return nil, apierror.NewFailedPreconditionError("scim is not enabled for the organization", fmt.Errorf("scim is not enabled for the organization"))
	}

	token := uuid.New()
	tokenSHA256 := sha256.Sum256(token[:])
	qSCIMAPIKey, err := q.CreateSCIMAPIKey(ctx, queries.CreateSCIMAPIKeyParams{
		ID:                uuid.New(),
		OrganizationID:    authn.OrganizationID(ctx),
		DisplayName:       req.ScimApiKey.DisplayName,
		SecretTokenSha256: tokenSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create scim api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	scimAPIKey := parseSCIMAPIKey(qSCIMAPIKey)
	scimAPIKey.SecretToken = idformat.SCIMAPIKeySecretToken.Format(token)
	return &frontendv1.CreateSCIMAPIKeyResponse{ScimApiKey: scimAPIKey}, nil
}

func (s *Store) UpdateSCIMAPIKey(ctx context.Context, req *frontendv1.UpdateSCIMAPIKeyRequest) (*frontendv1.UpdateSCIMAPIKeyResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	scimAPIKeyID, err := idformat.SCIMAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid scim api key id", fmt.Errorf("parse scim api key id: %w", err))
	}

	// authz
	qSCIMAPIKey, err := q.GetSCIMAPIKey(ctx, queries.GetSCIMAPIKeyParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             scimAPIKeyID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apierror.NewNotFoundError("scim api key not found", fmt.Errorf("get scim api key: %w", err))
		}
		return nil, fmt.Errorf("get scim api key: %w", err)
	}

	updates := queries.UpdateSCIMAPIKeyParams{
		ID:          scimAPIKeyID,
		DisplayName: qSCIMAPIKey.DisplayName,
	}

	if req.ScimApiKey.DisplayName != "" {
		updates.DisplayName = req.ScimApiKey.DisplayName
	}

	qUpdatedSCIMAPIKey, err := q.UpdateSCIMAPIKey(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update scim api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.UpdateSCIMAPIKeyResponse{ScimApiKey: parseSCIMAPIKey(qUpdatedSCIMAPIKey)}, nil
}

func (s *Store) DeleteSCIMAPIKey(ctx context.Context, req *frontendv1.DeleteSCIMAPIKeyRequest) (*frontendv1.DeleteSCIMAPIKeyResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	scimAPIKeyID, err := idformat.SCIMAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid scim api key id", fmt.Errorf("parse scim api key: %q", err))
	}

	// authz
	qSCIMAPIKey, err := q.GetSCIMAPIKey(ctx, queries.GetSCIMAPIKeyParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             scimAPIKeyID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apierror.NewNotFoundError("scim api key not found", fmt.Errorf("get scim api key: %w", err))
		}

		return nil, fmt.Errorf("get scim api key: %w", err)
	}

	if qSCIMAPIKey.SecretTokenSha256 != nil {
		return nil, apierror.NewFailedPreconditionError("scim api key must be revoked before deleting", fmt.Errorf("scim api key must be revoked before deleting"))
	}

	if err := q.DeleteSCIMAPIKey(ctx, scimAPIKeyID); err != nil {
		return nil, fmt.Errorf("delete scim api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.DeleteSCIMAPIKeyResponse{}, nil
}

func (s *Store) RevokeSCIMAPIKey(ctx context.Context, req *frontendv1.RevokeSCIMAPIKeyRequest) (*frontendv1.RevokeSCIMAPIKeyResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	scimAPIKeyID, err := idformat.SCIMAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid scim api key id", fmt.Errorf("parse scim api key id: %w", err))
	}

	// authz
	if _, err := q.GetSCIMAPIKey(ctx, queries.GetSCIMAPIKeyParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             scimAPIKeyID,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apierror.NewNotFoundError("scim api key not found", fmt.Errorf("get scim api key: %w", err))
		}

		return nil, fmt.Errorf("get scim api key: %w", err)
	}

	qSCIMAPIKey, err := q.RevokeSCIMAPIKey(ctx, scimAPIKeyID)
	if err != nil {
		return nil, fmt.Errorf("revoke scim api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.RevokeSCIMAPIKeyResponse{ScimApiKey: parseSCIMAPIKey(qSCIMAPIKey)}, nil
}

func parseSCIMAPIKey(qSCIMAPIKey queries.ScimApiKey) *frontendv1.SCIMAPIKey {
	return &frontendv1.SCIMAPIKey{
		Id:          idformat.SCIMAPIKey.Format(qSCIMAPIKey.ID),
		DisplayName: qSCIMAPIKey.DisplayName,
		CreateTime:  timestamppb.New(*qSCIMAPIKey.CreateTime),
		UpdateTime:  timestamppb.New(*qSCIMAPIKey.UpdateTime),
		SecretToken: "", // intentionally left blank
		Revoked:     qSCIMAPIKey.SecretTokenSha256 == nil,
	}
}
