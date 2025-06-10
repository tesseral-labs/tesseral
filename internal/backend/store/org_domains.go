package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) GetOrganizationDomains(ctx context.Context, req *backendv1.GetOrganizationDomainsRequest) (*backendv1.GetOrganizationDomainsResponse, error) {
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

	qDomains, err := q.GetOrganizationDomains(ctx, queries.GetOrganizationDomainsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	return &backendv1.GetOrganizationDomainsResponse{
		OrganizationDomains: parseOrganizationDomains(qOrg, qDomains),
	}, nil
}

func (s *Store) UpdateOrganizationDomains(ctx context.Context, req *backendv1.UpdateOrganizationDomainsRequest) (*backendv1.UpdateOrganizationDomainsResponse, error) {
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

	qPreviousDomains, err := q.GetOrganizationDomains(ctx, queries.GetOrganizationDomainsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	if err := q.DeleteOrganizationDomains(ctx, orgID); err != nil {
		return nil, fmt.Errorf("delete organization google hosted domains: %w", err)
	}

	for _, domain := range req.OrganizationDomains.Domains {
		if _, err := q.CreateOrganizationDomain(ctx, queries.CreateOrganizationDomainParams{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Domain:         domain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}

	qDomains, err := q.GetOrganizationDomains(ctx, queries.GetOrganizationDomainsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	organizationDomains := parseOrganizationDomains(qOrg, qDomains)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.organizations.update_domains",
		EventDetails: map[string]any{
			"domains":         organizationDomains.Domains,
			"previousDomains": parseOrganizationDomains(qOrg, qPreviousDomains).Domains,
		},
		OrganizationID: &qOrg.ID,
		ResourceType:   queries.AuditLogEventResourceTypeOrganization,
		ResourceID:     &qOrg.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateOrganizationDomainsResponse{
		OrganizationDomains: organizationDomains,
	}, nil
}

func parseOrganizationDomains(qOrg queries.Organization, qOrganizationDomains []queries.OrganizationDomain) *backendv1.OrganizationDomains {
	var Domains []string
	for _, qOrganizationDomain := range qOrganizationDomains {
		Domains = append(Domains, qOrganizationDomain.Domain)
	}
	return &backendv1.OrganizationDomains{
		OrganizationId: idformat.Organization.Format(qOrg.ID),
		Domains:        Domains,
	}
}
