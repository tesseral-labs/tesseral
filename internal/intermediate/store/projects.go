package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateProject(ctx context.Context, req *intermediatev1.CreateProjectRequest) (*intermediatev1.CreateProjectResponse, error) {
	if authn.ProjectID(ctx) != *s.dogfoodProjectID {
		return nil, apierror.NewPermissionDeniedError("cannot create a project", fmt.Errorf("create project attempted on non-dogfood project: %s", authn.ProjectID(ctx)))
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	if !intermediateSession.EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", nil)
	}

	// verify that the dogfood project exists
	qDogfoodProject, err := q.GetProjectByID(ctx, *s.dogfoodProjectID)
	if err != nil {
		return nil, fmt.Errorf("get dogfood project: %w", err)
	}

	// create this ahead of time so we can use it in the display name and auth domain
	newProjectID := uuid.New()
	formattedNewProjectID := idformat.Project.Format(newProjectID)
	newProjectVaultDomain := fmt.Sprintf("%s.%s", formattedNewProjectID, s.authAppsRootDomain)

	// create a new organization under the dogfood project
	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                 uuid.New(),
		DisplayName:        fmt.Sprintf("%s Backing Organization", formattedNewProjectID),
		ProjectID:          *s.dogfoodProjectID,
		LogInWithEmail:     qDogfoodProject.LogInWithEmail,
		LogInWithGoogle:    qDogfoodProject.LogInWithGoogle,
		LogInWithMicrosoft: qDogfoodProject.LogInWithMicrosoft,
		LogInWithPassword:  qDogfoodProject.LogInWithPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	// reflect the google hosted domain from the intermediate session if it exists
	if intermediateSession.GoogleHostedDomain != "" {
		if _, err := q.CreateOrganizationGoogleHostedDomain(ctx, queries.CreateOrganizationGoogleHostedDomainParams{
			OrganizationID:     qOrganization.ID,
			GoogleHostedDomain: intermediateSession.GoogleHostedDomain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}

	// reflect the microsoft tenant id from the intermediate session if it exists
	if intermediateSession.MicrosoftTenantId != "" {
		if _, err := q.CreateOrganizationMicrosoftTenantID(ctx, queries.CreateOrganizationMicrosoftTenantIDParams{
			OrganizationID:    qOrganization.ID,
			MicrosoftTenantID: intermediateSession.MicrosoftTenantId,
		}); err != nil {
			return nil, fmt.Errorf("create organization microsoft tenant id: %w", err)
		}
	}

	// create a new user invite for the intermediate session user
	if _, err := q.CreateUserInvite(ctx, queries.CreateUserInviteParams{
		ID:             uuid.New(),
		OrganizationID: qOrganization.ID,
		Email:          intermediateSession.Email,
		IsOwner:        true,
	}); err != nil {
		return nil, fmt.Errorf("create user invite: %w", err)
	}

	// create a new project backed by the new organization
	qProject, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                  newProjectID,
		OrganizationID:      &qOrganization.ID,
		VaultDomain:         newProjectVaultDomain,
		EmailSendFromDomain: fmt.Sprintf("mail.%s", s.authAppsRootDomain),
		DisplayName:         req.DisplayName,
		LogInWithGoogle:     false,
		LogInWithMicrosoft:  false,
		LogInWithPassword:   false,
		LogInWithSaml:       false,
	})
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	if err = commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.CreateProjectResponse{
		Project: parseProject(qProject),
	}, nil
}

func parseProject(qProject queries.Project) *intermediatev1.Project {
	return &intermediatev1.Project{
		Id:             qProject.ID.String(),
		OrganizationId: idformat.Organization.Format(*qProject.OrganizationID),
		CreateTime:     timestamppb.New(*qProject.CreateTime),
		UpdateTime:     timestamppb.New(*qProject.UpdateTime),
		DisplayName:    qProject.DisplayName,
		VaultDomain:    qProject.VaultDomain,
	}
}
