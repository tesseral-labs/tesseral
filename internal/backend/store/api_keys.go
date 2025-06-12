package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/prettysecret"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const apiKeySecretTokenSuffixLength = 4

func (s *Store) CreateAPIKey(ctx context.Context, req *backendv1.CreateAPIKeyRequest) (*backendv1.CreateAPIKeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.ApiKey.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qOrg, err := s.q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ID:        orgID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
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

	if !qOrg.ApiKeysEnabled {
		return nil, apierror.NewPermissionDeniedError("api keys are not enabled for this organization", fmt.Errorf("api keys not enabled for organization"))
	}

	var secretTokenValue [35]byte
	var secretToken string
	if _, err := rand.Read(secretTokenValue[:]); err != nil {
		return nil, fmt.Errorf("generate secret token: %w", err)
	}

	// Handle custom api key prefixes
	if qProject.ApiKeySecretTokenPrefix != nil {
		secretToken = prettysecret.Format(*qProject.ApiKeySecretTokenPrefix, secretTokenValue)
	}

	secretTokenSuffix := secretToken[len(secretToken)-apiKeySecretTokenSuffixLength:]
	secretTokenSHA256 := sha256.Sum256(secretTokenValue[:])

	var expireTime *time.Time
	if req.ApiKey.ExpireTime != nil {
		formattedExpireTime := req.ApiKey.ExpireTime.AsTime()
		expireTime = &formattedExpireTime
	}

	qAPIKey, err := q.CreateAPIKey(ctx, queries.CreateAPIKeyParams{
		ID:                uuid.New(),
		DisplayName:       req.ApiKey.DisplayName,
		ExpireTime:        expireTime,
		OrganizationID:    orgID,
		SecretTokenSha256: secretTokenSHA256[:],
		SecretTokenSuffix: &secretTokenSuffix,
	})
	if err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	apiKey := parseAPIKey(qAPIKey, &secretToken)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_keys.create",
		EventDetails: &backendv1.APIKeyCreated{
			ApiKey: apiKey,
		},
		OrganizationID: &qOrg.ID,
		ResourceType:   queries.AuditLogEventResourceTypeApiKey,
		ResourceID:     &qAPIKey.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.CreateAPIKeyResponse{
		ApiKey: apiKey,
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

	qApiKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
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
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("delete api key: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_keys.delete",
		EventDetails: &backendv1.APIKeyDeleted{
			ApiKey: parseAPIKey(qApiKey, nil),
		},
		OrganizationID: &qApiKey.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeApiKey,
		ResourceID:     &qApiKey.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
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
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	qAPIKeys, err := q.ListAPIKeys(ctx, queries.ListAPIKeysParams{
		ID:        orgID,
		ID_2:      startID,
		ProjectID: authn.ProjectID(ctx),
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}

	apiKeys := make([]*backendv1.APIKey, len(qAPIKeys))
	for i, qAPIKey := range qAPIKeys {
		apiKeys[i] = parseAPIKey(qAPIKey, nil)
	}

	var nextPageToken string
	if len(apiKeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qAPIKeys[limit].ID)
		apiKeys = apiKeys[:limit]
	}

	return &backendv1.ListAPIKeysResponse{
		ApiKeys:       apiKeys,
		NextPageToken: nextPageToken,
	}, nil
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

	qPreviousAPIKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("get api key by id: %w", err)
	}

	if err := q.RevokeAPIKey(ctx, queries.RevokeAPIKeyParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("revoke api key: %w", err)
	}

	qApiKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("get api key by id after revocation: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_keys.revoke",
		EventDetails: &backendv1.APIKeyRevoked{
			ApiKey:         parseAPIKey(qApiKey, nil),
			PreviousApiKey: parseAPIKey(qPreviousAPIKey, nil),
		},
		OrganizationID: &qApiKey.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeApiKey,
		ResourceID:     &qApiKey.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
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

	qPreviousApiKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get api key: %w", err))
		}
		return nil, fmt.Errorf("get api key by id: %w", err)
	}

	updatedApiKey, err := q.UpdateAPIKey(ctx, queries.UpdateAPIKeyParams{
		ID:          apiKeyID,
		DisplayName: req.ApiKey.DisplayName,
		ProjectID:   authn.ProjectID(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("update api key: %w", err)
	}

	apiKey := parseAPIKey(updatedApiKey, nil)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_keys.update",
		EventDetails: &backendv1.APIKeyUpdated{
			ApiKey:         apiKey,
			PreviousApiKey: parseAPIKey(qPreviousApiKey, nil),
		},
		OrganizationID: &updatedApiKey.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeApiKey,
		ResourceID:     &updatedApiKey.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.UpdateAPIKeyResponse{
		ApiKey: apiKey,
	}, nil
}

func (s *Store) AuthenticateAPIKey(ctx context.Context, req *backendv1.AuthenticateAPIKeyRequest) (*backendv1.AuthenticateAPIKeyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if !qProject.ApiKeysEnabled {
		return nil, apierror.NewPermissionDeniedError("api keys are not enabled for this project", fmt.Errorf("api keys not enabled for project"))
	}

	if qProject.ApiKeySecretTokenPrefix == nil {
		return nil, apierror.NewPermissionDeniedError("api key secret token prefix is not set for this project", fmt.Errorf("api key secret token prefix not set for project"))
	}

	secretTokenBytes, err := prettysecret.Parse(*qProject.ApiKeySecretTokenPrefix, req.SecretToken)
	if err != nil {
		return nil, apierror.NewUnauthenticatedApiKeyError("malformed_api_key_secret_token", fmt.Errorf("parse secret token: %w", err))
	}
	secretTokenSHA256 := sha256.Sum256(secretTokenBytes[:])

	qApiKeyDetails, err := q.GetAPIKeyDetailsBySecretTokenSHA256(ctx, queries.GetAPIKeyDetailsBySecretTokenSHA256Params{
		SecretTokenSha256: secretTokenSHA256[:],
		ProjectID:         authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewUnauthenticatedApiKeyError("invalid_api_key_secret_token", fmt.Errorf("get api key details: %w", err))
		}

		return nil, fmt.Errorf("get api key details: %w", err)
	}

	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ID:        qApiKeyDetails.OrganizationID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewUnauthenticatedApiKeyError("invalid_api_key_secret_token", fmt.Errorf("get organization: %w", err))
		}
		return nil, fmt.Errorf("get organization: %w", err)
	}

	if !qOrg.ApiKeysEnabled {
		return nil, apierror.NewPermissionDeniedError("api keys are not enabled for this organization", fmt.Errorf("api keys not enabled for organization"))
	}

	// Get all actions for the api key
	actions, err := q.GetAPIKeyActions(ctx, qApiKeyDetails.ID)
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	slices.Sort(actions)

	return &backendv1.AuthenticateAPIKeyResponse{
		ApiKeyId:       idformat.APIKey.Format(qApiKeyDetails.ID),
		Actions:        actions,
		OrganizationId: idformat.Organization.Format(qApiKeyDetails.OrganizationID),
	}, nil
}

func parseAPIKey(qAPIKey queries.ApiKey, secretToken *string) *backendv1.APIKey {
	return &backendv1.APIKey{
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
