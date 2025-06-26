package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
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

	qPreviousMicrosoftTenantIDs, err := q.GetOrganizationMicrosoftTenantIDs(ctx, queries.GetOrganizationMicrosoftTenantIDsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization microsoft tenant IDs: %w", err)
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

	microsoftTenantIDs := parseOrganizationMicrosoftTenantIDs(qOrg, qMicrosoftTenantIDs)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.organizations.update_microsoft_tenant_ids",
		EventDetails: &auditlogv1.UpdateOrganizationMicrosoftTenantIDs{
			MicrosoftTenantIds:         microsoftTenantIDs.MicrosoftTenantIds,
			PreviousMicrosoftTenantIds: parseOrganizationMicrosoftTenantIDs(qOrg, qPreviousMicrosoftTenantIDs).MicrosoftTenantIds,
		},
		OrganizationID: &qOrg.ID,
		ResourceType:   queries.AuditLogEventResourceTypeOrganization,
		ResourceID:     &qOrg.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateOrganizationMicrosoftTenantIDsResponse{
		OrganizationMicrosoftTenantIds: microsoftTenantIDs,
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
