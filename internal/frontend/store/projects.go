package store

import (
	"context"
	"fmt"

	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetProject(ctx context.Context, req *frontendv1.GetProjectRequest) (*frontendv1.GetProjectResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, projectid.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	return &frontendv1.GetProjectResponse{Project: parseProject(&qProject)}, nil
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
	}
}
