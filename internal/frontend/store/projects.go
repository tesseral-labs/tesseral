package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
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
	return &frontendv1.Project{
		Id:                        idformat.Project.Format(qProject.ID),
		CreateTime:                timestamppb.New(*qProject.CreateTime),
		UpdateTime:                timestamppb.New(*qProject.UpdateTime),
		DisplayName:               qProject.DisplayName,
		LogInWithGoogle:           qProject.LogInWithGoogle,
		LogInWithMicrosoft:        qProject.LogInWithMicrosoft,
		LogInWithGithub:           qProject.LogInWithGithub,
		LogInWithEmail:            qProject.LogInWithEmail,
		LogInWithPassword:         qProject.LogInWithPassword,
		LogInWithAuthenticatorApp: qProject.LogInWithAuthenticatorApp,
		LogInWithPasskey:          qProject.LogInWithPasskey,
		LogInWithSaml:             qProject.LogInWithSaml,
		VaultDomain:               qProject.VaultDomain,
		ApiKeysEnabled:            qProject.ApiKeysEnabled,
		ApiKeySecretTokenPrefix:   derefOrEmpty(qProject.ApiKeySecretTokenPrefix),
	}
}
