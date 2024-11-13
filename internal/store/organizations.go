package store

import (
	"context"

	"github.com/google/uuid"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

func (s *Store) CreateOrganization(ctx context.Context, req *openauthv1.Organization) (*openauthv1.Organization, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId, err := idformat.Organization.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	createdOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID: uuid.New(),
		ProjectID: projectId,
		DisplayName: req.DisplayName,
		GoogleHostedDomain: &req.GoogleHostedDomain,
		MicrosoftTenantID: &req.MicrosoftTenantId,
		OverrideLogInWithGoogleEnabled: &req.OverrideLogInWithGoogleEnabled,
		OverrideLogInWithMicrosoftEnabled: &req.OverrideLogInWithMicrosoftEnabled,
		OverrideLogInWithPasswordEnabled: &req.OverrideLogInWithPasswordEnabled,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformOrganization(createdOrganization), nil
}

func (s *Store) GetOrganization(ctx context.Context, req *openauthv1.ResourceIdRequest) (*openauthv1.Organization, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	organization, err := q.GetOrganizationByID(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	return transformOrganization(organization), nil
}

// TODO: Ensure that this function can only be called via a backend service reuqest
func (s *Store) ListOrganizations(ctx context.Context, req *openauthv1.ListOrganizationsRequest) (*openauthv1.ListOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId, err := idformat.Project.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	organizationRecords, err := q.ListOrganizationsByProjectId(ctx, queries.ListOrganizationsByProjectIdParams{
		ProjectID: projectId,
		Limit: int32(limit + 1),
	})
	if err != nil {
		return nil, err
	}

	organizations := []*openauthv1.Organization{}
	for _, organization := range organizationRecords {
		organizations = append(organizations, transformOrganization(organization))
	}

	var nextPageToken string
	if len(organizations) == limit + 1 {
		nextPageToken = s.pageEncoder.Marshal(organizations[limit].Id)
		organizations = organizations[:limit]
	}

	return &openauthv1.ListOrganizationsResponse{
		Organizations: organizations,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *openauthv1.Organization) (*openauthv1.Organization, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	existingOrganization, err := q.GetOrganizationByID(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	updates := queries.UpdateOrganizationParams{
		ID: organizationId,
	}

	// Conditionally update display name
	if req.DisplayName != "" {
		updates.DisplayName = req.DisplayName
	}

	// Conditionally update login method configs
	if req.GoogleHostedDomain != "" {
		updates.GoogleHostedDomain = &req.GoogleHostedDomain
	}
	if req.MicrosoftTenantId != "" {
		updates.MicrosoftTenantID = &req.MicrosoftTenantId
	}

	// Conditionally update overrides
	if req.OverrideLogInWithGoogleEnabled != *existingOrganization.OverrideLogInWithGoogleEnabled {
		updates.OverrideLogInWithGoogleEnabled = &req.OverrideLogInWithGoogleEnabled
	}
	if req.OverrideLogInWithMicrosoftEnabled != *existingOrganization.OverrideLogInWithMicrosoftEnabled {
		updates.OverrideLogInWithMicrosoftEnabled = &req.OverrideLogInWithMicrosoftEnabled
	}
	if req.OverrideLogInWithPasswordEnabled != *existingOrganization.OverrideLogInWithPasswordEnabled {
		updates.OverrideLogInWithPasswordEnabled = &req.OverrideLogInWithPasswordEnabled
	}

	updatedOrganization, err := q.UpdateOrganization(ctx, updates)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformOrganization(updatedOrganization), nil
}

func transformOrganization(organization queries.Organization) *openauthv1.Organization {
	return &openauthv1.Organization{
		Id: organization.ID.String(),
		ProjectId: organization.ProjectID.String(),
		DisplayName: organization.DisplayName,
		GoogleHostedDomain: *organization.GoogleHostedDomain,
		MicrosoftTenantId: *organization.MicrosoftTenantID,
		OverrideLogInWithGoogleEnabled: *organization.OverrideLogInWithGoogleEnabled,
		OverrideLogInWithMicrosoftEnabled: *organization.OverrideLogInWithMicrosoftEnabled,
		OverrideLogInWithPasswordEnabled: *organization.OverrideLogInWithPasswordEnabled,
	}
}