package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"golang.org/x/crypto/bcrypt"
)

func (s *Store) VerifyPassword(ctx context.Context, req *intermediatev1.VerifyPasswordRequest) (*intermediatev1.VerifyPasswordResponse, error) {
	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
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

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	// Ensure given organization is visible to intermediate session, and
	// suitable for authentication over password:
	//
	// 1. The organization must have passwords enabled,
	// 2. The intermediate session must have a verified email, and
	// 3. A user in that org must have the same email.
	passwordEnabled := qProject.LogInWithPasswordEnabled
	if qOrg.OverrideLogInMethods && qOrg.OverrideLogInWithPasswordEnabled != nil && !*qOrg.OverrideLogInWithPasswordEnabled {
		passwordEnabled = false
	}

	if !passwordEnabled {
		return nil, apierror.NewFailedPreconditionError("password authentication not enabled", nil)
	}

	emailVerified, err := s.getIntermediateSessionEmailVerified(ctx, q, qIntermediateSession.ID)
	if err != nil {
		return nil, fmt.Errorf("get intermediate session verified: %w", err)
	}

	if !emailVerified {
		return nil, apierror.NewFailedPreconditionError("email not verified", nil)
	}

	qMatchingUser, err := s.matchEmailUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match email user: %w", err)
	}

	if qMatchingUser == nil {
		return nil, apierror.NewFailedPreconditionError("no corresponding user found", nil)
	}

	if qMatchingUser.PasswordBcrypt == nil {
		return nil, apierror.NewFailedPreconditionError("user does not have password configured", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*qMatchingUser.PasswordBcrypt), []byte(req.Password)); err != nil {
		return nil, apierror.NewFailedPreconditionError("incorrect password", fmt.Errorf("bcrypt compare: %w", err))
	}

	if _, err := q.UpdateIntermediateSessionPasswordVerified(ctx, queries.UpdateIntermediateSessionPasswordVerifiedParams{
		OrganizationID: &qOrg.ID,
		ID:             qIntermediateSession.ID,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session password verified: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.VerifyPasswordResponse{}, nil
}
