package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/authn"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

func (s *Store) ListOrganizations(ctx context.Context, req *backendv1.ListOrganizationsRequest) (*backendv1.ListOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	organizationRecords, err := q.ListOrganizationsByProjectId(ctx, queries.ListOrganizationsByProjectIdParams{
		ProjectID: authn.ProjectID(ctx),
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}

	var organizations []*openauthv1.Organization
	for _, organization := range organizationRecords {
		organizations = append(organizations, parseOrganization(organization))
	}

	var nextPageToken string
	if len(organizations) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(organizations[limit].Id)
		organizations = organizations[:limit]
	}

	return &backendv1.ListOrganizationsResponse{
		Organizations: organizations,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetOrganization(ctx context.Context, req *backendv1.GetOrganizationRequest) (*openauthv1.Organization, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	organization, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        organizationId,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	return parseOrganization(organization), nil
}

func (s *Store) CreateOrganization(ctx context.Context, req *backendv1.CreateOrganizationRequest) (*openauthv1.Organization, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qOrg, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                                uuid.New(),
		ProjectID:                         authn.ProjectID(ctx),
		DisplayName:                       req.Organization.DisplayName,
		GoogleHostedDomain:                &req.Organization.GoogleHostedDomain,
		MicrosoftTenantID:                 &req.Organization.MicrosoftTenantId,
		OverrideLogInWithGoogleEnabled:    &req.Organization.OverrideLogInWithGoogleEnabled,
		OverrideLogInWithMicrosoftEnabled: &req.Organization.OverrideLogInWithMicrosoftEnabled,
		OverrideLogInWithPasswordEnabled:  &req.Organization.OverrideLogInWithPasswordEnabled,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseOrganization(qOrg), nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *backendv1.UpdateOrganizationRequest) (*openauthv1.Organization, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	// fetch existing org; this also acts as a permission check
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        organizationId,
	}); err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	updates := queries.UpdateOrganizationParams{
		ID: organizationId,
	}

	// Conditionally update display name
	if req.Organization.DisplayName != "" {
		updates.DisplayName = req.Organization.DisplayName
	}

	// TODO these don't work (nil deref), but I don't think they have the right
	// schema design anyway; skip for now
	//
	//// Conditionally update login method configs
	//if req.Organization.GoogleHostedDomain != "" {
	//	updates.GoogleHostedDomain = &req.Organization.GoogleHostedDomain
	//}
	//if req.Organization.MicrosoftTenantId != "" {
	//	updates.MicrosoftTenantID = &req.Organization.MicrosoftTenantId
	//}
	//
	//// Conditionally update overrides
	//if req.Organization.OverrideLogInWithGoogleEnabled != *existingOrganization.OverrideLogInWithGoogleEnabled {
	//	updates.OverrideLogInWithGoogleEnabled = &req.Organization.OverrideLogInWithGoogleEnabled
	//}
	//if req.Organization.OverrideLogInWithMicrosoftEnabled != *existingOrganization.OverrideLogInWithMicrosoftEnabled {
	//	updates.OverrideLogInWithMicrosoftEnabled = &req.Organization.OverrideLogInWithMicrosoftEnabled
	//}
	//if req.Organization.OverrideLogInWithPasswordEnabled != *existingOrganization.OverrideLogInWithPasswordEnabled {
	//	updates.OverrideLogInWithPasswordEnabled = &req.Organization.OverrideLogInWithPasswordEnabled
	//}

	updatedOrganization, err := q.UpdateOrganization(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseOrganization(updatedOrganization), nil
}

func parseOrganization(organization queries.Organization) *openauthv1.Organization {
	return &openauthv1.Organization{
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
