package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/emailaddr"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) ListOrganizations(ctx context.Context, req *intermediatev1.ListOrganizationsRequest) (*intermediatev1.ListOrganizationsResponse, error) {
	intermediateSession := authn.IntermediateSession(ctx)
	if !intermediateSession.EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", nil)
	}

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

	// Intermediate sessions can see an organization if:
	//
	// 1. The qOrg's google/microsoft hd/tid match, or
	// 2. There is a user in the qOrg with the same google/microsoft user id, or
	// 3. There is a user in the qOrg with the same email
	//
	// Options (2) and (3) are not redundant because a user may change their
	// email. The exchange endpoint will know to log the user in as the one that
	// has the same OAuth-based ID. It will also update that user's email
	// address.
	var qOrgs []queries.Organization

	if intermediateSession.GoogleHostedDomain != "" {
		// orgs with the same google hosted domain
		qGoogleOrgs, err := q.ListOrganizationsByGoogleHostedDomain(ctx, queries.ListOrganizationsByGoogleHostedDomainParams{
			ProjectID:          authn.ProjectID(ctx),
			GoogleHostedDomain: intermediateSession.GoogleHostedDomain,
		})
		if err != nil {
			return nil, fmt.Errorf("list organizations by google hosted domain: %w", err)
		}
		qOrgs = append(qOrgs, qGoogleOrgs...)
	}

	if intermediateSession.MicrosoftTenantId != "" {
		// orgs with the same microsoft tenant ID
		qMicrosoftOrgs, err := q.ListOrganizationsByMicrosoftTenantID(ctx, queries.ListOrganizationsByMicrosoftTenantIDParams{
			ProjectID:         authn.ProjectID(ctx),
			MicrosoftTenantID: intermediateSession.MicrosoftTenantId,
		})
		if err != nil {
			return nil, fmt.Errorf("list organizations by microsoft tenant id: %w", err)
		}
		qOrgs = append(qOrgs, qMicrosoftOrgs...)
	}

	// orgs with a matching user
	qUserOrgs, err := q.ListOrganizationsByMatchingUser(ctx, queries.ListOrganizationsByMatchingUserParams{
		ProjectID:       authn.ProjectID(ctx),
		Email:           intermediateSession.Email,
		GoogleUserID:    refOrNil(intermediateSession.GoogleUserId),
		MicrosoftUserID: refOrNil(intermediateSession.MicrosoftUserId),
	})
	if err != nil {
		return nil, fmt.Errorf("list organizations by matching user: %w", err)
	}
	qOrgs = append(qOrgs, qUserOrgs...)

	// dedupe qOrgs on ID
	var qOrgsDeduped []queries.Organization
	seen := map[uuid.UUID]struct{}{}
	for _, qOrg := range qOrgs {
		if _, ok := seen[qOrg.ID]; ok {
			continue
		}
		qOrgsDeduped = append(qOrgsDeduped, qOrg)
		seen[qOrg.ID] = struct{}{}
	}

	var organizations []*intermediatev1.Organization
	for _, organization := range qOrgsDeduped {
		qSamlConnection, err := q.GetOrganizationPrimarySAMLConnection(ctx, organization.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("primary saml connection not found", fmt.Errorf("get organization primary saml connection: %w", err))
		}

		organizations = append(organizations, parseOrganization(organization, qProject, &qSamlConnection))
	}

	return &intermediatev1.ListOrganizationsResponse{
		Organizations: organizations,
	}, nil
}

func (s *Store) ListSAMLOrganizations(ctx context.Context, req *intermediatev1.ListSAMLOrganizationsRequest) (*intermediatev1.ListSAMLOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	domain, err := emailaddr.Parse(req.Email)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid email address", fmt.Errorf("parse email: %w", err))
	}

	qOrganizations, err := q.ListSAMLOrganizations(ctx, queries.ListSAMLOrganizationsParams{
		ProjectID: authn.ProjectID(ctx),
		Domain:    domain,
	})
	if err != nil {
		return nil, fmt.Errorf("list saml organizations: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	var organizations []*intermediatev1.Organization
	for _, organization := range qOrganizations {
		qSamlConnection, err := q.GetOrganizationPrimarySAMLConnection(ctx, organization.ID)
		if err != nil {
			return nil, fmt.Errorf("get organization primary saml connection: %w", err)
		}

		organizations = append(organizations, parseOrganization(organization, qProject, &qSamlConnection))
	}

	return &intermediatev1.ListSAMLOrganizationsResponse{
		Organizations: organizations,
	}, nil
}

func parseOrganization(qOrg queries.Organization, qProject queries.Project, qSAMLConnection *queries.SamlConnection) *intermediatev1.Organization {
	logInWithGoogleEnabled := qProject.LogInWithGoogleEnabled
	logInWithMicrosoftEnabled := qProject.LogInWithMicrosoftEnabled
	logInWithPasswordEnabled := qProject.LogInWithPasswordEnabled

	// allow orgs to disable login methods
	if derefOrEmpty(qOrg.DisableLogInWithGoogle) {
		logInWithGoogleEnabled = false
	}
	if derefOrEmpty(qOrg.DisableLogInWithMicrosoft) {
		logInWithMicrosoftEnabled = false
	}
	if derefOrEmpty(qOrg.DisableLogInWithPassword) {
		logInWithPasswordEnabled = false
	}

	return &intermediatev1.Organization{
		Id:                        idformat.Organization.Format(qOrg.ID),
		DisplayName:               qOrg.DisplayName,
		LogInWithGoogleEnabled:    logInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: logInWithMicrosoftEnabled,
		LogInWithPasswordEnabled:  logInWithPasswordEnabled,
		PrimarySamlConnectionId:   idformat.SAMLConnection.Format(qSAMLConnection.ID),
	}
}
