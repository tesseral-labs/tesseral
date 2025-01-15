package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	return &frontendv1.GetProjectResponse{Project: parseProject(&qProject)}, nil
}

func parseProject(qProject *queries.Project) *frontendv1.Project {
	authDomain := derefOrEmpty(qProject.AuthDomain)
	if qProject.CustomAuthDomain != nil {
		authDomain = *qProject.CustomAuthDomain
	}

	return &frontendv1.Project{
		Id:                        idformat.Project.Format(qProject.ID),
		CreateTime:                timestamppb.New(*qProject.CreateTime),
		UpdateTime:                timestamppb.New(*qProject.UpdateTime),
		DisplayName:               qProject.DisplayName,
		LogInWithPasswordEnabled:  qProject.LogInWithPasswordEnabled,
		LogInWithGoogleEnabled:    qProject.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: qProject.LogInWithMicrosoftEnabled,
		AuthDomain:                authDomain,
	}
}
