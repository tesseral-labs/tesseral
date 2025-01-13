package store

import (
	"context"

	"github.com/google/uuid"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetProjectIDByDomain(ctx context.Context, domain string) (*uuid.UUID, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := q.GetProjectIDByCustomDomain(ctx, []string{domain})
	if err != nil {
		return nil, err
	}

	return &projectID, nil
}

func parseProject(qProject *queries.Project) *intermediatev1.Project {
	return &intermediatev1.Project{
		Id:                        idformat.Project.Format(qProject.ID),
		LogInWithPasswordEnabled:  qProject.LogInWithPasswordEnabled,
		LogInWithGoogleEnabled:    qProject.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: qProject.LogInWithMicrosoftEnabled,
		CustomDomains:             qProject.CustomDomains,
	}
}
