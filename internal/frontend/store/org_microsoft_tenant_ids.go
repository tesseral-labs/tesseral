package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
)

func (s *Store) GetOrganizationMicrosoftTenantIDs(ctx context.Context, req *frontendv1.GetOrganizationMicrosoftTenantIDsRequest) (*frontendv1.GetOrganizationMicrosoftTenantIDsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qMicrosoftTenantIDs, err := q.GetOrganizationMicrosoftTenantIDs(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization microsoft tenant ids: %w", err)
	}

	return &frontendv1.GetOrganizationMicrosoftTenantIDsResponse{
		OrganizationMicrosoftTenantIds: parseOrganizationMicrosoftTenantIDs(qMicrosoftTenantIDs),
	}, nil
}

func (s *Store) UpdateOrganizationMicrosoftTenantIDs(ctx context.Context, req *frontendv1.UpdateOrganizationMicrosoftTenantIDsRequest) (*frontendv1.UpdateOrganizationMicrosoftTenantIDsResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.DeleteOrganizationMicrosoftTenantIDs(ctx, authn.OrganizationID(ctx)); err != nil {
		return nil, fmt.Errorf("delete organization microsoft tenant ids: %w", err)
	}

	for _, microsoftTenantID := range req.OrganizationMicrosoftTenantIds.MicrosoftTenantIds {
		if _, err := q.CreateOrganizationMicrosoftTenantID(ctx, queries.CreateOrganizationMicrosoftTenantIDParams{
			ID:                uuid.New(),
			OrganizationID:    authn.OrganizationID(ctx),
			MicrosoftTenantID: microsoftTenantID,
		}); err != nil {
			return nil, fmt.Errorf("create organization microsoft tenant id: %w", err)
		}
	}

	qMicrosoftTenantIDs, err := q.GetOrganizationMicrosoftTenantIDs(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization microsoft tenant ids: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.UpdateOrganizationMicrosoftTenantIDsResponse{
		OrganizationMicrosoftTenantIds: parseOrganizationMicrosoftTenantIDs(qMicrosoftTenantIDs),
	}, nil
}

func parseOrganizationMicrosoftTenantIDs(qOrganizationMicrosoftTenantIDs []queries.OrganizationMicrosoftTenantID) *frontendv1.OrganizationMicrosoftTenantIDs {
	var microsoftTenantIDs []string
	for _, qOrganizationMicrosoftTenantID := range qOrganizationMicrosoftTenantIDs {
		microsoftTenantIDs = append(microsoftTenantIDs, qOrganizationMicrosoftTenantID.MicrosoftTenantID)
	}
	return &frontendv1.OrganizationMicrosoftTenantIDs{
		MicrosoftTenantIds: microsoftTenantIDs,
	}
}
