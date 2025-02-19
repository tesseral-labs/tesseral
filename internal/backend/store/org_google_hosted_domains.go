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

func (s *Store) GetOrganizationGoogleHostedDomains(ctx context.Context, req *backendv1.GetOrganizationGoogleHostedDomainsRequest) (*backendv1.GetOrganizationGoogleHostedDomainsResponse, error) {
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

	qGoogleHostedDomains, err := q.GetOrganizationGoogleHostedDomains(ctx, queries.GetOrganizationGoogleHostedDomainsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	return &backendv1.GetOrganizationGoogleHostedDomainsResponse{
		OrganizationGoogleHostedDomains: parseOrganizationGoogleHostedDomains(qOrg, qGoogleHostedDomains),
	}, nil
}

func (s *Store) UpdateOrganizationGoogleHostedDomains(ctx context.Context, req *backendv1.UpdateOrganizationGoogleHostedDomainsRequest) (*backendv1.UpdateOrganizationGoogleHostedDomainsResponse, error) {
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

	if err := q.DeleteOrganizationGoogleHostedDomains(ctx, orgID); err != nil {
		return nil, fmt.Errorf("delete organization google hosted domains: %w", err)
	}

	for _, googleHostedDomain := range req.OrganizationGoogleHostedDomains.GoogleHostedDomains {
		if _, err := q.CreateOrganizationGoogleHostedDomain(ctx, queries.CreateOrganizationGoogleHostedDomainParams{
			ID:                 uuid.New(),
			OrganizationID:     orgID,
			GoogleHostedDomain: googleHostedDomain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}

	qGoogleHostedDomains, err := q.GetOrganizationGoogleHostedDomains(ctx, queries.GetOrganizationGoogleHostedDomainsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization google hosted domains: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateOrganizationGoogleHostedDomainsResponse{
		OrganizationGoogleHostedDomains: parseOrganizationGoogleHostedDomains(qOrg, qGoogleHostedDomains),
	}, nil
}

func parseOrganizationGoogleHostedDomains(qOrg queries.Organization, qOrganizationGoogleHostedDomains []queries.OrganizationGoogleHostedDomain) *backendv1.OrganizationGoogleHostedDomains {
	var googleHostedDomains []string
	for _, qOrganizationGoogleHostedDomain := range qOrganizationGoogleHostedDomains {
		googleHostedDomains = append(googleHostedDomains, qOrganizationGoogleHostedDomain.GoogleHostedDomain)
	}
	return &backendv1.OrganizationGoogleHostedDomains{
		OrganizationId:      idformat.Organization.Format(qOrg.ID),
		GoogleHostedDomains: googleHostedDomains,
	}
}
