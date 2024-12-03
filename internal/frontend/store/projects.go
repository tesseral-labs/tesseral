package store

import (
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func parseProject(project *queries.Project) *frontendv1.Project {
	return &frontendv1.Project{
		Id:                         idformat.Project.Format(project.ID),
		OrganizationId:             idformat.Organization.Format(*project.OrganizationID),
		LogInWithPasswordEnabled:   project.LogInWithPasswordEnabled,
		LogInWithGoogleEnabled:     project.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled:  project.LogInWithMicrosoftEnabled,
		GoogleOauthClientId:        derefOrEmpty(project.GoogleOauthClientID),
		GoogleOauthClientSecret:    derefOrEmpty(project.GoogleOauthClientSecret),
		MicrosoftOauthClientId:     derefOrEmpty(project.MicrosoftOauthClientID),
		MicrosoftOauthClientSecret: derefOrEmpty(project.MicrosoftOauthClientSecret),
	}
}
