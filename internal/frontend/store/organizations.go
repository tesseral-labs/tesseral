package store

import (
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func parseOrganization(organization queries.Organization) *frontendv1.Organization {
	return &frontendv1.Organization{
		Id:                                idformat.Organization.Format(organization.ID),
		ProjectId:                         idformat.Project.Format(organization.ProjectID),
		DisplayName:                       organization.DisplayName,
		GoogleHostedDomain:                derefOrEmpty(organization.GoogleHostedDomain),
		MicrosoftTenantId:                 derefOrEmpty(organization.MicrosoftTenantID),
		OverrideLogInWithGoogleEnabled:    derefOrEmpty(organization.OverrideLogInWithGoogleEnabled),
		OverrideLogInWithMicrosoftEnabled: derefOrEmpty(organization.OverrideLogInWithMicrosoftEnabled),
		OverrideLogInWithPasswordEnabled:  derefOrEmpty(organization.OverrideLogInWithPasswordEnabled),
	}
}
