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

func (s *Store) CreateOrganization(ctx context.Context, req *intermediatev1.CreateOrganizationRequest) (*intermediatev1.CreateOrganizationResponse, error) {
	intermediateSession := authn.IntermediateSession(ctx)
	intermediateSessionID, err := idformat.IntermediateSession.Parse(intermediateSession.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid intermediate session ID", fmt.Errorf("parse intermediate session ID: %w", err))
	}

	if !intermediateSession.EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", nil)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                 uuid.New(),
		ProjectID:          authn.ProjectID(ctx),
		DisplayName:        req.DisplayName,
		LogInWithGoogle:    qProject.LogInWithGoogle,
		LogInWithMicrosoft: qProject.LogInWithMicrosoft,
		LogInWithPassword:  qProject.LogInWithPassword,
		ScimEnabled:        false,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	// If the intermediate session is associated with a Google or Microsoft
	// login, associate the organization as well.
	if intermediateSession.GoogleHostedDomain != "" {
		if _, err := q.CreateOrganizationGoogleHostedDomain(ctx, queries.CreateOrganizationGoogleHostedDomainParams{
			ID:                 uuid.New(),
			OrganizationID:     qOrganization.ID,
			GoogleHostedDomain: intermediateSession.GoogleHostedDomain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}
	if intermediateSession.MicrosoftTenantId != "" {
		if _, err := q.CreateOrganizationMicrosoftTenantID(ctx, queries.CreateOrganizationMicrosoftTenantIDParams{
			ID:                uuid.New(),
			OrganizationID:    qOrganization.ID,
			MicrosoftTenantID: intermediateSession.MicrosoftTenantId,
		}); err != nil {
			return nil, fmt.Errorf("create organization microsoft tenant id: %w", err)
		}
	}

	if _, err = q.UpdateIntermediateSessionOrganizationID(ctx, queries.UpdateIntermediateSessionOrganizationIDParams{
		ID:             intermediateSessionID,
		OrganizationID: &qOrganization.ID,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session organization ID: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.CreateOrganizationResponse{}, nil
}

func (s *Store) ListOrganizations(ctx context.Context, req *intermediatev1.ListOrganizationsRequest) (*intermediatev1.ListOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("intermediate session not found", fmt.Errorf("get intermediate session by id: %w", err))
		}

		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	if !authn.IntermediateSession(ctx).EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", fmt.Errorf("intermediate session has unverified email"))
	}

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

	if qIntermediateSession.GoogleHostedDomain != nil {
		// orgs with the same google hosted domain
		qGoogleOrgs, err := q.ListOrganizationsByGoogleHostedDomain(ctx, queries.ListOrganizationsByGoogleHostedDomainParams{
			ProjectID:          authn.ProjectID(ctx),
			GoogleHostedDomain: *qIntermediateSession.GoogleHostedDomain,
		})
		if err != nil {
			return nil, fmt.Errorf("list organizations by google hosted domain: %w", err)
		}
		qOrgs = append(qOrgs, qGoogleOrgs...)
	}

	if qIntermediateSession.MicrosoftTenantID != nil {
		// orgs with the same microsoft tenant ID
		qMicrosoftOrgs, err := q.ListOrganizationsByMicrosoftTenantID(ctx, queries.ListOrganizationsByMicrosoftTenantIDParams{
			ProjectID:         authn.ProjectID(ctx),
			MicrosoftTenantID: *qIntermediateSession.MicrosoftTenantID,
		})
		if err != nil {
			return nil, fmt.Errorf("list organizations by microsoft tenant id: %w", err)
		}
		qOrgs = append(qOrgs, qMicrosoftOrgs...)
	}

	// orgs with a matching user
	qUserOrgs, err := q.ListOrganizationsByMatchingUser(ctx, queries.ListOrganizationsByMatchingUserParams{
		ProjectID:       authn.ProjectID(ctx),
		Email:           *qIntermediateSession.Email,
		GoogleUserID:    qIntermediateSession.GoogleUserID,
		MicrosoftUserID: qIntermediateSession.MicrosoftUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("list organizations by matching user: %w", err)
	}
	qOrgs = append(qOrgs, qUserOrgs...)

	// orgs with a matching user invite
	qUserInviteOrgs, err := q.ListOrganizationsByMatchingUserInvite(ctx, queries.ListOrganizationsByMatchingUserInviteParams{
		ProjectID: authn.ProjectID(ctx),
		Email:     *qIntermediateSession.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("list organizations by matching user invite: %w", err)
	}
	qOrgs = append(qOrgs, qUserInviteOrgs...)

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
	for _, qOrg := range qOrgsDeduped {
		qSamlConnection, err := q.GetOrganizationPrimarySAMLConnection(ctx, qOrg.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("primary saml connection not found", fmt.Errorf("get organization primary saml connection: %w", err))
		}

		samlConnectionOrNil := &qSamlConnection
		if qSamlConnection.ID == uuid.Nil {
			samlConnectionOrNil = nil
		}

		// Parse the organization before performing additional checks
		org := parseOrganization(qOrg, qProject, samlConnectionOrNil)

		// Check if the user exists on the organization.
		existingUser, err := s.matchEmailUser(ctx, q, qOrg, qIntermediateSession)
		if err != nil {
			return nil, fmt.Errorf("match email user: %w", err)
		}

		org.UserExists = existingUser != nil
		if existingUser != nil {
			org.UserHasPassword = existingUser.PasswordBcrypt != nil
			org.UserHasAuthenticatorApp = existingUser.AuthenticatorAppSecretCiphertext != nil

			hasPasskeys, err := q.GetUserHasPasskey(ctx, existingUser.ID)
			if err != nil {
				return nil, fmt.Errorf("get user has passkey: %w", err)
			}

			org.UserHasPasskey = hasPasskeys
		}

		// Append the parsed organization to the list of organizations.
		organizations = append(organizations, org)
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
	for _, qOrg := range qOrganizations {
		qSamlConnection, err := q.GetOrganizationPrimarySAMLConnection(ctx, qOrg.ID)
		if err != nil {
			return nil, fmt.Errorf("get organization primary saml connection: %w", err)
		}

		// Append the parsed organization to the list of organizations.
		organizations = append(organizations, parseOrganization(qOrg, qProject, &qSamlConnection))
	}

	return &intermediatev1.ListSAMLOrganizationsResponse{
		Organizations: organizations,
	}, nil
}

func (s *Store) SetOrganization(ctx context.Context, req *intermediatev1.SetOrganizationRequest) (*intermediatev1.SetOrganizationResponse, error) {
	intermediateSession := authn.IntermediateSession(ctx)
	intermediateSessionID, err := idformat.IntermediateSession.Parse(intermediateSession.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid intermediate session ID", fmt.Errorf("parse intermediate session ID: %w", err))
	}

	if intermediateSession.OrganizationId != "" {
		return nil, apierror.NewFailedPreconditionError("organization already set", fmt.Errorf("organization already set"))
	}

	if !intermediateSession.EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", nil)
	}

	organizationID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization ID", fmt.Errorf("parse organization ID: %w", err))
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qOrganization, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ID:        organizationID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	if _, err = q.UpdateIntermediateSessionOrganizationID(ctx, queries.UpdateIntermediateSessionOrganizationIDParams{
		ID:             intermediateSessionID,
		OrganizationID: &qOrganization.ID,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session organization ID: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.SetOrganizationResponse{}, nil
}

func parseOrganization(qOrg queries.Organization, qProject queries.Project, qSAMLConnection *queries.SamlConnection) *intermediatev1.Organization {
	var primarySamlConnectionID string
	if qSAMLConnection != nil {
		primarySamlConnectionID = idformat.SAMLConnection.Format(qSAMLConnection.ID)
	}

	return &intermediatev1.Organization{
		Id:                        idformat.Organization.Format(qOrg.ID),
		DisplayName:               qOrg.DisplayName,
		LogInWithGoogle:           qOrg.LogInWithGoogle,
		LogInWithMicrosoft:        qOrg.LogInWithMicrosoft,
		LogInWithPassword:         qOrg.LogInWithPassword,
		LogInWithAuthenticatorApp: qOrg.LogInWithAuthenticatorApp,
		LogInWithPasskey:          qOrg.LogInWithPasskey,
		RequireMfa:                qOrg.RequireMfa,
		PrimarySamlConnectionId:   primarySamlConnectionID,
	}
}
