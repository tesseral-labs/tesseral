package store

import (
	"context"
	"fmt"
	"time"

	"github.com/openauth/openauth/internal/bcryptcost"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"golang.org/x/crypto/bcrypt"
)

const (
	// after this many failed attempts, lock out a user
	passwordLockoutAttempts = 5
	// how long to lock users out
	passwordLockoutDuration = time.Minute * 10
)

func (s *Store) RegisterPassword(ctx context.Context, req *intermediatev1.RegisterPasswordRequest) (*intermediatev1.RegisterPasswordResponse, error) {
	intermediateSession := authn.IntermediateSession(ctx)
	if intermediateSession.OrganizationId != "" {
		return nil, apierror.NewFailedPreconditionError("organization id already set for intermediate session", fmt.Errorf("organization id already set for intermediate session"))
	}

	if intermediateSession.PasswordVerified {
		return nil, apierror.NewFailedPreconditionError("user already verified for intermediate session", fmt.Errorf("user already verified for intermediate session"))
	}

	orgID, err := idformat.Organization.Parse(intermediateSession.OrganizationId)
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

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	if err := enforceOrganizationLoginEnabled(qOrg); err != nil {
		return nil, fmt.Errorf("enforce organization login enabled: %w", err)
	}

	// Ensure given organization is suitable for authentication over password:
	passwordEnabled := qProject.LogInWithPasswordEnabled
	if qProject.LogInWithPasswordEnabled || (qOrg.OverrideLogInMethods && qOrg.OverrideLogInWithPasswordEnabled != nil && !*qOrg.OverrideLogInWithPasswordEnabled) {
		passwordEnabled = false
	}
	if !passwordEnabled {
		return nil, apierror.NewFailedPreconditionError("password authentication not enabled", fmt.Errorf("password authentication not enabled"))
	}

	emailVerified, err := s.getIntermediateSessionEmailVerified(ctx, q, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session email verified: %w", err)
	}
	if !emailVerified {
		return nil, apierror.NewFailedPreconditionError("email not verified", fmt.Errorf("email not verified"))
	}

	passwordBcryptBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptcost.Cost)
	if err != nil {
		return nil, fmt.Errorf("generate bcrypt hash: %w", err)
	}

	passwordBcrypt := string(passwordBcryptBytes)
	_, err = q.UpdateIntermediateSessionNewUserPasswordBcrypt(ctx, queries.UpdateIntermediateSessionNewUserPasswordBcryptParams{
		ID:                    authn.IntermediateSessionID(ctx),
		NewUserPasswordBcrypt: &passwordBcrypt,
	})
	if err != nil {
		return nil, fmt.Errorf("update intermediate session new user password bcrypt: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RegisterPasswordResponse{}, nil
}

func (s *Store) VerifyPassword(ctx context.Context, req *intermediatev1.VerifyPasswordRequest) (*intermediatev1.VerifyPasswordResponse, error) {
	intermediateSession := authn.IntermediateSession(ctx)
	if intermediateSession.OrganizationId != "" {
		return nil, apierror.NewFailedPreconditionError("organization id already set for intermediate session", fmt.Errorf("organization id already set for intermediate session"))
	}

	if intermediateSession.PasswordVerified {
		return nil, apierror.NewFailedPreconditionError("user already verified for intermediate session", fmt.Errorf("user already verified for intermediate session"))
	}

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

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	if err := enforceOrganizationLoginEnabled(qOrg); err != nil {
		return nil, fmt.Errorf("enforce organization login enabled: %w", err)
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
	if derefOrEmpty(qOrg.DisableLogInWithPassword) {
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

	if qMatchingUser.PasswordLockoutExpireTime != nil && qMatchingUser.PasswordLockoutExpireTime.After(time.Now()) {
		return nil, apierror.NewFailedPreconditionError("too many password attempts; user is temporarily locked out", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*qMatchingUser.PasswordBcrypt), []byte(req.Password)); err != nil {
		attempts := qMatchingUser.FailedPasswordAttempts + 1
		if attempts >= passwordLockoutAttempts {
			// lock the user out
			passwordLockoutExpireTime := time.Now().Add(passwordLockoutDuration)
			if _, err := q.UpdateUserPasswordLockoutExpireTime(ctx, queries.UpdateUserPasswordLockoutExpireTimeParams{
				ID:                        qMatchingUser.ID,
				PasswordLockoutExpireTime: &passwordLockoutExpireTime,
			}); err != nil {
				return nil, err
			}

			// reset fail count
			if _, err := q.UpdateUserFailedPasswordAttempts(ctx, queries.UpdateUserFailedPasswordAttemptsParams{
				ID:                     qMatchingUser.ID,
				FailedPasswordAttempts: 0,
			}); err != nil {
				return nil, err
			}

			if err := commit(); err != nil {
				return nil, fmt.Errorf("commit: %w", err)
			}

			return nil, apierror.NewFailedPreconditionError("too many password attempts; user is temporarily locked out", nil)
		}

		// update fail count, but do not lock out
		if _, err := q.UpdateUserFailedPasswordAttempts(ctx, queries.UpdateUserFailedPasswordAttemptsParams{
			ID:                     qMatchingUser.ID,
			FailedPasswordAttempts: attempts,
		}); err != nil {
			return nil, fmt.Errorf("update user failed password attempts: %w", err)
		}

		if err := commit(); err != nil {
			return nil, fmt.Errorf("commit: %w", err)
		}

		return nil, apierror.NewFailedPreconditionError("incorrect password", fmt.Errorf("bcrypt compare: %w", err))
	}

	// Re-write password back to database; this lets us progressively increase
	// bcrypt costs over time.
	//
	// We could avoid these writes by checking the PasswordBcrypt using
	// bcrypt.Cost, but for relatively small additional cost, not doing so
	// reduces the complexity and number of paths through this code.
	passwordBcryptBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptcost.Cost)
	if err != nil {
		return nil, fmt.Errorf("generate bcrypt hash: %w", err)
	}

	passwordBcrypt := string(passwordBcryptBytes)
	if _, err := q.UpdateUserPasswordBcrypt(ctx, queries.UpdateUserPasswordBcryptParams{
		ID:             qMatchingUser.ID,
		PasswordBcrypt: &passwordBcrypt,
	}); err != nil {
		return nil, fmt.Errorf("update user password bcrypt: %w", err)
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
