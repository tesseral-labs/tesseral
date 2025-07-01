package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) ConsoleSearch(ctx context.Context, req *backendv1.ConsoleSearchRequest) (*backendv1.ConsoleSearchResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	var (
		apiKeyID         *uuid.UUID
		backendAPIKeyID  *uuid.UUID
		organizationID   *uuid.UUID
		publishableKeyID *uuid.UUID
		userID           *uuid.UUID
	)

	parsedAPIKeyID, err := idformat.APIKey.Parse(req.Query)
	if err == nil {
		apiKeyID = (*uuid.UUID)(&parsedAPIKeyID)
	}
	parsedBackendAPIKeyID, err := idformat.BackendAPIKey.Parse(req.Query)
	if err == nil {
		backendAPIKeyID = (*uuid.UUID)(&parsedBackendAPIKeyID)
	}
	parsedOrganizationID, err := idformat.Organization.Parse(req.Query)
	if err == nil {
		organizationID = (*uuid.UUID)(&parsedOrganizationID)
	}
	parsedPublishableKeyID, err := idformat.PublishableKey.Parse(req.Query)
	if err == nil {
		publishableKeyID = (*uuid.UUID)(&parsedPublishableKeyID)
	}
	parsedUserID, err := idformat.User.Parse(req.Query)
	if err == nil {
		userID = (*uuid.UUID)(&parsedUserID)
	}

	limit := 5
	if req.Limit > 0 {
		limit = int(req.Limit)
	}

	qAPIKeys, err := q.ConsoleSearchAPIKeys(ctx, queries.ConsoleSearchAPIKeysParams{
		ID:        apiKeyID,
		Limit:     int32(limit),
		ProjectID: authn.ProjectID(ctx),
		Query:     req.Query,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("search API keys: %w", err)
		}
	}
	var apiKeys []*backendv1.APIKey
	for _, qAPIKey := range qAPIKeys {
		apiKeys = append(apiKeys, parseAPIKey(qAPIKey))
	}

	qBackendAPIKeys, err := q.ConsoleSearchBackendAPIKeys(ctx, queries.ConsoleSearchBackendAPIKeysParams{
		ID:        backendAPIKeyID,
		Limit:     int32(limit),
		ProjectID: authn.ProjectID(ctx),
		Query:     req.Query,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("search backend API keys: %w", err)
		}
	}
	var backendAPIKeys []*backendv1.BackendAPIKey
	for _, backendAPIKey := range qBackendAPIKeys {
		backendAPIKeys = append(backendAPIKeys, parseBackendAPIKey(backendAPIKey))
	}

	qOrganizations, err := q.ConsoleSearchOrganizations(ctx, queries.ConsoleSearchOrganizationsParams{
		ID:        organizationID,
		Limit:     int32(limit),
		ProjectID: authn.ProjectID(ctx),
		Query:     req.Query,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("search organizations: %w", err)
		}
	}
	var organizations []*backendv1.Organization
	for _, qOrganization := range qOrganizations {
		organizations = append(organizations, parseOrganization(qProject, qOrganization))
	}

	qPublishableKeys, err := q.ConsoleSearchPublishableKeys(ctx, queries.ConsoleSearchPublishableKeysParams{
		ID:        publishableKeyID,
		Limit:     int32(limit),
		ProjectID: authn.ProjectID(ctx),
		Query:     req.Query,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("search publishable keys: %w", err)
		}
	}
	var publishableKeys []*backendv1.PublishableKey
	for _, qPublishableKey := range qPublishableKeys {
		publishableKeys = append(publishableKeys, parsePublishableKey(qPublishableKey))
	}

	qUsers, err := q.ConsoleSearchUsers(ctx, queries.ConsoleSearchUsersParams{
		ID:        userID,
		Limit:     int32(limit),
		ProjectID: authn.ProjectID(ctx),
		Query:     req.Query,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("search users: %w", err)
		}
	}
	var users []*backendv1.User
	for _, qUser := range qUsers {
		users = append(users, parseUser(qUser))
	}

	return &backendv1.ConsoleSearchResponse{
		ApiKeys:         apiKeys,
		BackendApiKeys:  backendAPIKeys,
		Organizations:   organizations,
		PublishableKeys: publishableKeys,
		Users:           users,
	}, nil
}
