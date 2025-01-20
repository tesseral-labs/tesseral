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

func (s *Store) GetOrganization(ctx context.Context, req *frontendv1.GetOrganizationRequest) (*frontendv1.GetOrganizationResponse, error) {
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

	qOrganization, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
	}

	return &frontendv1.GetOrganizationResponse{Organization: parseOrganization(qProject, qOrganization)}, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *frontendv1.UpdateOrganizationRequest) (*frontendv1.UpdateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	qOrg, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	updates := queries.UpdateOrganizationParams{
		ID: authn.OrganizationID(ctx),
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

	qUpdatedOrg, err := q.UpdateOrganization(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", fmt.Errorf("update organization: %w", err))
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	// Commit the transaction
	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.UpdateOrganizationResponse{
		Organization: parseOrganization(qProject, qUpdatedOrg),
	}, nil
}

func parseOrganization(qProject queries.Project, qOrg queries.Organization) *frontendv1.Organization {
	logInWithGoogleEnabled := qProject.LogInWithGoogleEnabled
	logInWithMicrosoftEnabled := qProject.LogInWithMicrosoftEnabled
	logInWithPasswordEnabled := qProject.LogInWithPasswordEnabled

	if qOrg.OverrideLogInMethods {
		logInWithGoogleEnabled = derefOrEmpty(qOrg.OverrideLogInWithGoogleEnabled)
		logInWithMicrosoftEnabled = derefOrEmpty(qOrg.OverrideLogInWithMicrosoftEnabled)
		logInWithPasswordEnabled = derefOrEmpty(qOrg.OverrideLogInWithPasswordEnabled)
	}

	return &frontendv1.Organization{
		Id:                        idformat.Organization.Format(qOrg.ID),
		DisplayName:               qOrg.DisplayName,
		CreateTime:                timestamppb.New(*qOrg.CreateTime),
		UpdateTime:                timestamppb.New(*qOrg.UpdateTime),
		OverrideLogInMethods:      &qOrg.OverrideLogInMethods,
		LogInWithGoogleEnabled:    logInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: logInWithMicrosoftEnabled,
		LogInWithPasswordEnabled:  logInWithPasswordEnabled,
		// TODO these are a list now
		//GoogleHostedDomain:        derefOrEmpty(qOrg.GoogleHostedDomain),
		//MicrosoftTenantId:         derefOrEmpty(qOrg.MicrosoftTenantID),
		SamlEnabled: qOrg.SamlEnabled,
		ScimEnabled: qOrg.ScimEnabled,
	}
}

// validateIsOwner returns an error if the current user is not an owner of the
// organization.
func (s *Store) validateIsOwner(ctx context.Context) error {
	qUser, err := s.q.GetUserByID(ctx, authn.UserID(ctx))
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	if !qUser.IsOwner {
		return apierror.NewPermissionDeniedError("user must be an owner of the organization", fmt.Errorf("user is not an owner"))
	}
	return nil
}
