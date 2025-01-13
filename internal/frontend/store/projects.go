package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetProject(ctx context.Context, req *frontendv1.GetProjectRequest) (*frontendv1.GetProjectResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	return &frontendv1.GetProjectResponse{Project: parseProject(&qProject)}, nil
}

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

func parseProject(qProject *queries.Project) *frontendv1.Project {
	return &frontendv1.Project{
		Id:                        idformat.Project.Format(qProject.ID),
		CreateTime:                timestamppb.New(*qProject.CreateTime),
		UpdateTime:                timestamppb.New(*qProject.UpdateTime),
		DisplayName:               qProject.DisplayName,
		LogInWithPasswordEnabled:  qProject.LogInWithPasswordEnabled,
		LogInWithGoogleEnabled:    qProject.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: qProject.LogInWithMicrosoftEnabled,
		CustomDomains:             qProject.CustomDomains,
	}
}
