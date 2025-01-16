package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateOrganization(ctx context.Context, req *backendv1.CreateOrganizationRequest) (*backendv1.CreateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var googleHostedDomain *string
	if req.Organization.GoogleHostedDomain != "" {
		googleHostedDomain = &req.Organization.GoogleHostedDomain
	}

	var microsoftTenantId *string
	if req.Organization.MicrosoftTenantId != "" {
		microsoftTenantId = &req.Organization.MicrosoftTenantId
	}

	var (
		overrideLogInWithGoogleEnabled,
		overrideLogInWithMicrosoftEnabled,
		overrideLogInWithPasswordEnabled *bool
	)

	if req.Organization.OverrideLogInMethods != nil {
		overrideLogInWithGoogleEnabled = &req.Organization.LogInWithGoogleEnabled
		overrideLogInWithMicrosoftEnabled = &req.Organization.LogInWithMicrosoftEnabled
		overrideLogInWithPasswordEnabled = &req.Organization.LogInWithPasswordEnabled
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	samlEnabled := qProject.OrganizationsSamlEnabledDefault
	if req.Organization.SamlEnabled != nil {
		samlEnabled = *req.Organization.SamlEnabled
	}

	scimEnabled := qProject.OrganizationsScimEnabledDefault
	if req.Organization.ScimEnabled != nil {
		scimEnabled = *req.Organization.ScimEnabled
	}

	qOrg, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                 uuid.New(),
		ProjectID:          authn.ProjectID(ctx),
		DisplayName:        req.Organization.DisplayName,
		GoogleHostedDomain: googleHostedDomain,
		MicrosoftTenantID:  microsoftTenantId,

		OverrideLogInMethods:              derefOrEmpty(req.Organization.OverrideLogInMethods),
		OverrideLogInWithGoogleEnabled:    overrideLogInWithGoogleEnabled,
		OverrideLogInWithMicrosoftEnabled: overrideLogInWithMicrosoftEnabled,
		OverrideLogInWithPasswordEnabled:  overrideLogInWithPasswordEnabled,

		SamlEnabled: samlEnabled,
		ScimEnabled: scimEnabled,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.CreateOrganizationResponse{Organization: parseOrganization(qProject, qOrg)}, nil
}

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

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	limit := 10
	qOrgs, err := q.ListOrganizationsByProjectId(ctx, queries.ListOrganizationsByProjectIdParams{
		ProjectID: authn.ProjectID(ctx),
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}

	var organizations []*backendv1.Organization
	for _, qOrg := range qOrgs {
		organizations = append(organizations, parseOrganization(qProject, qOrg))
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

func (s *Store) GetOrganization(ctx context.Context, req *backendv1.GetOrganizationRequest) (*backendv1.GetOrganizationResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        organizationId,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	return &backendv1.GetOrganizationResponse{Organization: parseOrganization(qProject, qOrg)}, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *backendv1.UpdateOrganizationRequest) (*backendv1.UpdateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// fetch existing org; this also acts as a permission check
	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	updates := queries.UpdateOrganizationParams{
		ID: orgID,
	}

	updates.DisplayName = qOrg.DisplayName
	if req.Organization.DisplayName != "" {
		updates.DisplayName = req.Organization.DisplayName
	}

	updates.OverrideLogInMethods = qOrg.OverrideLogInMethods
	if req.Organization.OverrideLogInMethods != nil {
		updates.OverrideLogInMethods = *req.Organization.OverrideLogInMethods
	}

	// update the override_log_in_with_..._enabled columns to null unless the
	// organization is overriding those columns.
	if req.Organization.GetOverrideLogInMethods() {
		updates.OverrideLogInWithGoogleEnabled = &req.Organization.LogInWithGoogleEnabled
		updates.OverrideLogInWithMicrosoftEnabled = &req.Organization.LogInWithMicrosoftEnabled
		updates.OverrideLogInWithPasswordEnabled = &req.Organization.LogInWithPasswordEnabled
	}

	updates.SamlEnabled = qOrg.SamlEnabled
	if req.Organization.SamlEnabled != nil {
		updates.SamlEnabled = *req.Organization.SamlEnabled
	}

	updates.ScimEnabled = qOrg.ScimEnabled
	if req.Organization.ScimEnabled != nil {
		updates.ScimEnabled = *req.Organization.ScimEnabled
	}

	qUpdatedOrg, err := q.UpdateOrganization(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateOrganizationResponse{Organization: parseOrganization(qProject, qUpdatedOrg)}, nil
}

func (s *Store) DeleteOrganization(ctx context.Context, req *backendv1.DeleteOrganizationRequest) (*backendv1.DeleteOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// authz check
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	if err := q.DeleteOrganization(ctx, orgID); err != nil {
		return nil, fmt.Errorf("delete organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteOrganizationResponse{}, nil
}

func parseOrganization(qProject queries.Project, qOrg queries.Organization) *backendv1.Organization {
	logInWithGoogleEnabled := qProject.LogInWithGoogleEnabled
	logInWithMicrosoftEnabled := qProject.LogInWithMicrosoftEnabled
	logInWithPasswordEnabled := qProject.LogInWithPasswordEnabled

	if qOrg.OverrideLogInMethods {
		// Only allow overrides to restrict settings, not augment them. We can't
		// easily enforce this rule in UpdateOrganization, because a project may
		// retroactively remove support for a login method. Such an update
		// should immediately affect all organizations.
		logInWithGoogleEnabled = qProject.LogInWithGoogleEnabled && derefOrEmpty(qOrg.OverrideLogInWithGoogleEnabled)
		logInWithMicrosoftEnabled = qProject.LogInWithMicrosoftEnabled && derefOrEmpty(qOrg.OverrideLogInWithMicrosoftEnabled)
		logInWithPasswordEnabled = qProject.LogInWithPasswordEnabled && derefOrEmpty(qOrg.OverrideLogInWithPasswordEnabled)
	}

	return &backendv1.Organization{
		Id:                        idformat.Organization.Format(qOrg.ID),
		DisplayName:               qOrg.DisplayName,
		CreateTime:                timestamppb.New(*qOrg.CreateTime),
		UpdateTime:                timestamppb.New(*qOrg.UpdateTime),
		OverrideLogInMethods:      &qOrg.OverrideLogInMethods,
		LogInWithPasswordEnabled:  logInWithPasswordEnabled,
		LogInWithGoogleEnabled:    logInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: logInWithMicrosoftEnabled,
		GoogleHostedDomain:        derefOrEmpty(qOrg.GoogleHostedDomain),
		MicrosoftTenantId:         derefOrEmpty(qOrg.MicrosoftTenantID),
		SamlEnabled:               &qOrg.SamlEnabled,
		ScimEnabled:               &qOrg.ScimEnabled,
	}
}
