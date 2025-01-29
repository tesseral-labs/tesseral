package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetOrganizationMicrosoftTenantIDs(ctx context.Context, req *backendv1.GetOrganizationMicrosoftTenantIDsRequest) (*backendv1.GetOrganizationMicrosoftTenantIDsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	qMicrosoftTenantIDs, err := q.GetOrganizationMicrosoftTenantIDs(ctx, queries.GetOrganizationMicrosoftTenantIDsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization microsoft tenant IDs: %w", err)
	}

	var microsoftTenantIDs []string
	for _, qMicrosoftTenantID := range qMicrosoftTenantIDs {
		microsoftTenantIDs = append(microsoftTenantIDs, qMicrosoftTenantID.MicrosoftTenantID)
	}

	return &backendv1.GetOrganizationMicrosoftTenantIDsResponse{
		OrganizationMicrosoftTenantIds: parseOrganizationMicrosoftTenantIDs(qOrg, qMicrosoftTenantIDs),
	}, nil
}

func (s *Store) UpdateOrganizationMicrosoftTenantIDs(ctx context.Context, req *backendv1.UpdateOrganizationMicrosoftTenantIDsRequest) (*backendv1.UpdateOrganizationMicrosoftTenantIDsResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	if err := q.DeleteOrganizationMicrosoftTenantIDs(ctx, orgID); err != nil {
		return nil, fmt.Errorf("delete organization microsoft tenant IDs: %w", err)
	}

	for _, microsoftTenantID := range req.OrganizationMicrosoftTenantIds.MicrosoftTenantIds {
		if _, err := q.CreateOrganizationMicrosoftTenantID(ctx, queries.CreateOrganizationMicrosoftTenantIDParams{
			ID:                uuid.New(),
			OrganizationID:    orgID,
			MicrosoftTenantID: microsoftTenantID,
		}); err != nil {
			return nil, fmt.Errorf("create organization microsoft tenant ID: %w", err)
		}
	}

	qMicrosoftTenantIDs, err := q.GetOrganizationMicrosoftTenantIDs(ctx, queries.GetOrganizationMicrosoftTenantIDsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization microsoft tenant IDs: %w", err)
	}

	var microsoftTenantIDs []string
	for _, qMicrosoftTenantID := range qMicrosoftTenantIDs {
		microsoftTenantIDs = append(microsoftTenantIDs, qMicrosoftTenantID.MicrosoftTenantID)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateOrganizationMicrosoftTenantIDsResponse{
		OrganizationMicrosoftTenantIds: parseOrganizationMicrosoftTenantIDs(qOrg, qMicrosoftTenantIDs),
	}, nil
}

func parseOrganizationMicrosoftTenantIDs(qOrg queries.Organization, qOrganizationMicrosoftTenantIDs []queries.OrganizationMicrosoftTenantID) *backendv1.OrganizationMicrosoftTenantIDs {
	var microsoftTenantIDs []string
	for _, qOrganizationMicrosoftTenantID := range qOrganizationMicrosoftTenantIDs {
		microsoftTenantIDs = append(microsoftTenantIDs, qOrganizationMicrosoftTenantID.MicrosoftTenantID)
	}
	return &backendv1.OrganizationMicrosoftTenantIDs{
		OrganizationId:     idformat.Organization.Format(qOrg.ID),
		MicrosoftTenantIds: microsoftTenantIDs,
	}
}
