package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetOrganization(ctx context.Context, req *frontendv1.GetOrganizationRequest) (*frontendv1.GetOrganizationResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, projectid.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrganization, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	return &frontendv1.GetOrganizationResponse{Organization: parseOrganization(qProject, qOrganization)}, nil
}

func parseOrganization(qProject queries.Project, qOrg queries.Organization) *frontendv1.Organization {
	logInWithGoogleEnabled := qProject.LogInWithGoogleEnabled
	logInWithMicrosoftEnabled := qProject.LogInWithMicrosoftEnabled
	logInWithPasswordEnabled := qProject.LogInWithPasswordEnabled

	if qOrg.OverrideLogInMethods {
		logInWithGoogleEnabled = derefOrEmpty(qOrg.OverrideLogInWithGoogleEnabled)
		logInWithMicrosoftEnabled = derefOrEmpty(qOrg.OverrideLogInWithMicrosoftEnabled)
		logInWithPasswordEnabled = derefOrEmpty(qOrg.OverrideLogInWithPasswordEnabled)
	}

	return &frontendv1.Organization{
		Id:                        idformat.Organization.Format(qOrg.ID),
		ProjectId:                 idformat.Project.Format(qOrg.ProjectID),
		DisplayName:               qOrg.DisplayName,
		GoogleHostedDomain:        derefOrEmpty(qOrg.GoogleHostedDomain),
		MicrosoftTenantId:         derefOrEmpty(qOrg.MicrosoftTenantID),
		OverrideLogInMethods:      &qOrg.OverrideLogInMethods,
		LogInWithGoogleEnabled:    logInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: logInWithMicrosoftEnabled,
		LogInWithPasswordEnabled:  logInWithPasswordEnabled,
	}
}
