package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/muststructpb"
)

func (s *Store) GetOrganizationGoogleHostedDomains(ctx context.Context, req *frontendv1.GetOrganizationGoogleHostedDomainsRequest) (*frontendv1.GetOrganizationGoogleHostedDomainsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qGoogleHostedDomains, err := q.GetOrganizationGoogleHostedDomains(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	return &frontendv1.GetOrganizationGoogleHostedDomainsResponse{
		OrganizationGoogleHostedDomains: parseOrganizationGoogleHostedDomains(qGoogleHostedDomains),
	}, nil
}

func (s *Store) UpdateOrganizationGoogleHostedDomains(ctx context.Context, req *frontendv1.UpdateOrganizationGoogleHostedDomainsRequest) (*frontendv1.UpdateOrganizationGoogleHostedDomainsResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// Get the current organization google hosted domains before deleting them to log the changes
	qPreviousGoogleHostedDomains, err := q.GetOrganizationGoogleHostedDomains(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	if err := q.DeleteOrganizationGoogleHostedDomains(ctx, authn.OrganizationID(ctx)); err != nil {
		return nil, fmt.Errorf("delete organization google hosted domains: %w", err)
	}

	for _, googleHostedDomain := range req.OrganizationGoogleHostedDomains.GoogleHostedDomains {
		if _, err := q.CreateOrganizationGoogleHostedDomain(ctx, queries.CreateOrganizationGoogleHostedDomainParams{
			ID:                 uuid.New(),
			OrganizationID:     authn.OrganizationID(ctx),
			GoogleHostedDomain: googleHostedDomain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}

	qGoogleHostedDomains, err := q.GetOrganizationGoogleHostedDomains(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	googleHostedDomains := parseOrganizationGoogleHostedDomains(qGoogleHostedDomains)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.organizations.update_google_hosted_domains",
		EventDetails: muststructpb.MustNewValue(map[string]any{
			"googleHostedDomains":         googleHostedDomains.GoogleHostedDomains,
			"previousGoogleHostedDomains": parseOrganizationGoogleHostedDomains(qPreviousGoogleHostedDomains).GoogleHostedDomains,
		}),
		ResourceType: queries.AuditLogEventResourceTypeOrganization,
		ResourceID:   refOrNil(authn.OrganizationID(ctx)),
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.UpdateOrganizationGoogleHostedDomainsResponse{
		OrganizationGoogleHostedDomains: googleHostedDomains,
	}, nil
}

func parseOrganizationGoogleHostedDomains(qOrganizationGoogleHostedDomains []queries.OrganizationGoogleHostedDomain) *frontendv1.OrganizationGoogleHostedDomains {
	var googleHostedDomains []string
	for _, qOrganizationGoogleHostedDomain := range qOrganizationGoogleHostedDomains {
		googleHostedDomains = append(googleHostedDomains, qOrganizationGoogleHostedDomain.GoogleHostedDomain)
	}
	return &frontendv1.OrganizationGoogleHostedDomains{
		GoogleHostedDomains: googleHostedDomains,
	}
}
