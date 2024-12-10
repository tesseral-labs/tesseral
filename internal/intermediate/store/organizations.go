package store

import (
	"context"

	"github.com/google/uuid"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) CreateOrganization(ctx context.Context, req *intermediatev1.CreateOrganizationRequest) (*intermediatev1.CreateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId := projectid.ProjectID(ctx)
	project, err := q.GetProjectByID(ctx, projectId)
	if err != nil {
		return nil, err
	}

	createdOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                                uuid.New(),
		ProjectID:                         projectId,
		DisplayName:                       req.DisplayName,
		OverrideLogInWithGoogleEnabled:    &project.LogInWithGoogleEnabled,
		OverrideLogInWithMicrosoftEnabled: &project.LogInWithMicrosoftEnabled,
		OverrideLogInWithPasswordEnabled:  &project.LogInWithPasswordEnabled,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.CreateOrganizationResponse{
		Organization: parseOrganization(createdOrganization),
	}, nil
}

func (s *Store) ListOrganizations(
	ctx context.Context,
	req *intermediatev1.ListOrganizationsRequest,
) (*intermediatev1.ListOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId, err := idformat.Project.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	limit := 10
	organizationRecords, err := q.ListOrganizationsByProjectIdAndEmail(ctx, queries.ListOrganizationsByProjectIdAndEmailParams{
		ProjectID: projectId,
		Email:     req.Email,
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, err
	}

	organizations := []*intermediatev1.Organization{}
	for _, organization := range organizationRecords {
		organizations = append(organizations, parseOrganization(organization))
	}

	var nextPageToken string
	if len(organizations) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(organizations[limit].Id)
		organizations = organizations[:limit]
	}

	return &intermediatev1.ListOrganizationsResponse{
		Organizations: organizations,
		NextPageToken: nextPageToken,
	}, nil
}

func parseOrganization(organization queries.Organization) *intermediatev1.Organization {
	return &intermediatev1.Organization{
		Id:                        organization.ID.String(),
		DisplayName:               organization.DisplayName,
		LogInWithGoogleEnabled:    *organization.OverrideLogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: *organization.OverrideLogInWithMicrosoftEnabled,
		LogInWithPasswordEnabled:  *organization.OverrideLogInWithPasswordEnabled,
	}
}
