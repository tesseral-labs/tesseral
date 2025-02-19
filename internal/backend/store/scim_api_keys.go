package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListSCIMAPIKeys(ctx context.Context, req *backendv1.ListSCIMAPIKeysRequest) (*backendv1.ListSCIMAPIKeysResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// authz
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qSCIMAPIKeys, err := q.ListSCIMAPIKeys(ctx, queries.ListSCIMAPIKeysParams{
		OrganizationID: orgID,
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list scim api keys: %w", err)
	}

	var scimAPIKeys []*backendv1.SCIMAPIKey
	for _, qSCIMAPIKey := range qSCIMAPIKeys {
		scimAPIKeys = append(scimAPIKeys, parseSCIMAPIKey(qSCIMAPIKey))
	}

	var nextPageToken string
	if len(scimAPIKeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qSCIMAPIKeys[limit].ID)
		scimAPIKeys = scimAPIKeys[:limit]
	}

	return &backendv1.ListSCIMAPIKeysResponse{
		ScimApiKeys:   scimAPIKeys,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetSCIMAPIKey(ctx context.Context, req *backendv1.GetSCIMAPIKeyRequest) (*backendv1.GetSCIMAPIKeyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	scimAPIKeyID, err := idformat.SCIMAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid scim api key id", fmt.Errorf("parse scim api key id: %w", err))
	}

	qSCIMAPIKey, err := q.GetSCIMAPIKey(ctx, queries.GetSCIMAPIKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        scimAPIKeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("scim api key not found", fmt.Errorf("get scim api key: %w", err))
		}

		return nil, fmt.Errorf("get scim api key: %w", err)
	}

	return &backendv1.GetSCIMAPIKeyResponse{ScimApiKey: parseSCIMAPIKey(qSCIMAPIKey)}, nil
}

func (s *Store) CreateSCIMAPIKey(ctx context.Context, req *backendv1.CreateSCIMAPIKeyRequest) (*backendv1.CreateSCIMAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.ScimApiKey.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// authz
	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	if !qOrg.ScimEnabled {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("organization does not have SCIM enabled"))
	}

	token := uuid.New()
	tokenSHA256 := sha256.Sum256(token[:])
	qSCIMAPIKey, err := q.CreateSCIMAPIKey(ctx, queries.CreateSCIMAPIKeyParams{
		ID:                uuid.New(),
		OrganizationID:    orgID,
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
	return &backendv1.CreateSCIMAPIKeyResponse{ScimApiKey: scimAPIKey}, nil
}

func (s *Store) UpdateSCIMAPIKey(ctx context.Context, req *backendv1.UpdateSCIMAPIKeyRequest) (*backendv1.UpdateSCIMAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	scimAPIKeyID, err := idformat.SCIMAPIKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid scim api key", fmt.Errorf("parse scim api key id: %w", err))
	}

	// authz
	qSCIMAPIKey, err := q.GetSCIMAPIKey(ctx, queries.GetSCIMAPIKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        scimAPIKeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

	return &backendv1.UpdateSCIMAPIKeyResponse{ScimApiKey: parseSCIMAPIKey(qUpdatedSCIMAPIKey)}, nil
}

func (s *Store) DeleteSCIMAPIKey(ctx context.Context, req *backendv1.DeleteSCIMAPIKeyRequest) (*backendv1.DeleteSCIMAPIKeyResponse, error) {
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
		ProjectID: authn.ProjectID(ctx),
		ID:        scimAPIKeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("scim api key not found", fmt.Errorf("get scim api key: %w", err))
		}

		return nil, fmt.Errorf("get scim api key: %w", err)
	}

	if qSCIMAPIKey.SecretTokenSha256 != nil {
		return nil, apierror.NewFailedPreconditionError("scim api key must be revoked before deletion", fmt.Errorf("scim api key must be revoked before deletion"))
	}

	if err := q.DeleteSCIMAPIKey(ctx, scimAPIKeyID); err != nil {
		return nil, fmt.Errorf("delete scim api key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteSCIMAPIKeyResponse{}, nil
}

func (s *Store) RevokeSCIMAPIKey(ctx context.Context, req *backendv1.RevokeSCIMAPIKeyRequest) (*backendv1.RevokeSCIMAPIKeyResponse, error) {
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
		ProjectID: authn.ProjectID(ctx),
		ID:        scimAPIKeyID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

	return &backendv1.RevokeSCIMAPIKeyResponse{ScimApiKey: parseSCIMAPIKey(qSCIMAPIKey)}, nil
}

func parseSCIMAPIKey(qSCIMAPIKey queries.ScimApiKey) *backendv1.SCIMAPIKey {
	return &backendv1.SCIMAPIKey{
		Id:             idformat.SCIMAPIKey.Format(qSCIMAPIKey.ID),
		OrganizationId: idformat.Organization.Format(qSCIMAPIKey.OrganizationID),
		DisplayName:    qSCIMAPIKey.DisplayName,
		CreateTime:     timestamppb.New(*qSCIMAPIKey.CreateTime),
		UpdateTime:     timestamppb.New(*qSCIMAPIKey.UpdateTime),
		SecretToken:    "", // intentionally left blank
		Revoked:        qSCIMAPIKey.SecretTokenSha256 == nil,
	}
}
