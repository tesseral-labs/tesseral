package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

const sessionDuration = time.Hour * 24 * 7

func (s *Store) ExchangeIntermediateSessionForSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        *qIntermediateSession.OrganizationID,
	})
	if err != nil {
		return nil, err
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	if err := enforceOrganizationLoginEnabled(qOrg); err != nil {
		return nil, fmt.Errorf("enforce organization login enabled: %w", err)
	}

	if err := s.validateAuthRequirementsSatisfied(ctx, q, qIntermediateSession.ID); err != nil {
		return nil, fmt.Errorf("validate auth requirements satisfied: %w", err)
	}

	qUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match user: %w", err)
	}

	// if no matching user, create a new one
	if qUser == nil {
		qNewUser, err := q.CreateUser(ctx, queries.CreateUserParams{
			ID:              uuid.New(),
			OrganizationID:  qOrg.ID,
			Email:           *qIntermediateSession.Email,
			GoogleUserID:    qIntermediateSession.GoogleUserID,
			MicrosoftUserID: qIntermediateSession.MicrosoftUserID,
			PasswordBcrypt:  qIntermediateSession.NewUserPasswordBcrypt,
		})
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}

		qUser = &qNewUser
	}

	// if a passkey is registered on the intermediate session, copy it onto the
	// user
	if qIntermediateSession.PasskeyCredentialID != nil {
		if err := s.copyRegisteredPasskeySettings(ctx, q, qIntermediateSession, *qUser); err != nil {
			return nil, fmt.Errorf("copy registered passkey settings: %w", err)
		}
	}

	// if an authenticator app is registered on the intermediate session, copy
	// it onto the user
	if qIntermediateSession.AuthenticatorAppSecretCiphertext != nil {
		if err := s.copyRegisteredAuthenticatorAppSettings(ctx, q, qIntermediateSession, *qUser); err != nil {
			return nil, fmt.Errorf("copy registered authenticator app settings: %w", err)
		}
	}

	expireTime := time.Now().Add(sessionDuration)

	// Create a new session for the user
	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])
	if _, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expireTime,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qUser.ID,
	}); err != nil {
		return nil, err
	}

	// revoke the intermediate session
	if _, err := q.RevokeIntermediateSession(ctx, qIntermediateSession.ID); err != nil {
		return nil, fmt.Errorf("revoke intermediate session: %w", err)
	}

	// delete any outstanding invites for this email
	qUserInvite, err := q.DeleteIntermediateSessionUserInvite(ctx, queries.DeleteIntermediateSessionUserInviteParams{
		OrganizationID: *qIntermediateSession.OrganizationID,
		Email:          *qIntermediateSession.Email,
	})
	if err != nil {
		// this error is ok; it just means there's no user invite
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("delete intermediate session user invite: %w", err)
		}
	}

	if qUserInvite.ID != uuid.Nil {
		if qUserInvite.IsOwner {
			if _, err := q.UpdateUserIsOwner(ctx, queries.UpdateUserIsOwnerParams{
				ID:      qUser.ID,
				IsOwner: true,
			}); err != nil {
				return nil, fmt.Errorf("update user is owner: %w", err)
			}
		}
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.ExchangeIntermediateSessionForSessionResponse{
		AccessToken:  "", // populated in service
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func (s *Store) validateAuthRequirementsSatisfied(ctx context.Context, q *queries.Queries, intermediateSessionID uuid.UUID) error {
	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, intermediateSessionID)
	if err != nil {
		return fmt.Errorf("get intermediate session by id: %w", err)
	}

	if qIntermediateSession.OrganizationID == nil {
		return apierror.NewFailedPreconditionError("organization not set", fmt.Errorf("organization not set"))
	}

	emailVerified, err := s.getIntermediateSessionEmailVerified(ctx, q, qIntermediateSession.ID)
	if err != nil {
		return fmt.Errorf("get intermediate session verified: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        *qIntermediateSession.OrganizationID,
	})
	if err != nil {
		return fmt.Errorf("get organization by id: %w", err)
	}

	return validateAuthRequirementsSatisfiedInner(qIntermediateSession, emailVerified, qOrg)
}

func validateAuthRequirementsSatisfiedInner(qIntermediateSession queries.IntermediateSession, emailVerified bool, qOrg queries.Organization) error {
	if qIntermediateSession.Email == nil {
		panic(fmt.Errorf("intermediate session missing email: %v", qIntermediateSession.ID))
	}

	if qIntermediateSession.PrimaryLoginFactor == nil {
		return apierror.NewFailedPreconditionError("primary login factor not set", nil)
	}

	if !emailVerified {
		return apierror.NewFailedPreconditionError("email not verified", nil)
	}

	if qOrg.LogInWithPassword && !qIntermediateSession.PasswordVerified {
		return apierror.NewFailedPreconditionError("password not verified", nil)
	}

	if qOrg.RequireMfa {
		hasPasskey := qOrg.LogInWithPasskey && qIntermediateSession.PasskeyVerified
		hasAuthenticatorApp := qOrg.LogInWithAuthenticatorApp && qIntermediateSession.AuthenticatorAppVerified

		if !hasPasskey && !hasAuthenticatorApp {
			return apierror.NewFailedPreconditionError("mfa required", nil)
		}
	}

	switch *qIntermediateSession.PrimaryLoginFactor {
	case queries.PrimaryLoginFactorGoogleOauth:
		if qIntermediateSession.GoogleUserID == nil {
			panic(fmt.Errorf("intermediate session missing google user id: %v", qIntermediateSession.ID))
		}

		if qOrg.LogInWithGoogle {
			return nil
		}
	case queries.PrimaryLoginFactorMicrosoftOauth:
		if qIntermediateSession.MicrosoftUserID == nil {
			panic(fmt.Errorf("intermediate session missing microsoft user id: %v", qIntermediateSession.ID))
		}

		if qOrg.LogInWithMicrosoft {
			return nil
		}
	case queries.PrimaryLoginFactorEmail:
		if qOrg.LogInWithEmail {
			return nil
		}
	}

	return apierror.NewFailedPreconditionError("no authentication method satisfied", nil)
}

func (s *Store) matchUser(ctx context.Context, q *queries.Queries, qOrg queries.Organization, qIntermediateSession queries.IntermediateSession) (*queries.User, error) {
	qUser, err := s.matchGoogleUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match google user: %w", err)
	}
	if qUser != nil {
		return qUser, nil
	}

	qUser, err = s.matchMicrosoftUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match microsoft user: %w", err)
	}
	if qUser != nil {
		return qUser, nil
	}

	qUser, err = s.matchEmailUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match email user: %w", err)
	}
	if qUser != nil {
		return qUser, nil
	}

	return nil, nil
}

func (s *Store) matchGoogleUser(ctx context.Context, q *queries.Queries, qOrg queries.Organization, qIntermediateSession queries.IntermediateSession) (*queries.User, error) {
	if qIntermediateSession.GoogleUserID == nil {
		return nil, nil
	}

	qUser, err := q.GetOrganizationUserByGoogleUserID(ctx, queries.GetOrganizationUserByGoogleUserIDParams{
		OrganizationID: qOrg.ID,
		GoogleUserID:   qIntermediateSession.GoogleUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get organization user by google user id: %w", err)
	}

	return &qUser, nil
}

func (s *Store) matchMicrosoftUser(ctx context.Context, q *queries.Queries, qOrg queries.Organization, qIntermediateSession queries.IntermediateSession) (*queries.User, error) {
	if qIntermediateSession.MicrosoftUserID == nil {
		return nil, nil
	}

	qUser, err := q.GetOrganizationUserByMicrosoftUserID(ctx, queries.GetOrganizationUserByMicrosoftUserIDParams{
		OrganizationID:  qOrg.ID,
		MicrosoftUserID: qIntermediateSession.MicrosoftUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get organization user by microsoft user id: %w", err)
	}

	return &qUser, nil
}

func (s *Store) matchEmailUser(ctx context.Context, q *queries.Queries, qOrg queries.Organization, qIntermediateSession queries.IntermediateSession) (*queries.User, error) {
	qUser, err := q.GetOrganizationUserByEmail(ctx, queries.GetOrganizationUserByEmailParams{
		OrganizationID: qOrg.ID,
		Email:          *qIntermediateSession.Email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get organization user by email user id: %w", err)
	}

	return &qUser, nil
}

func (s *Store) copyRegisteredPasskeySettings(ctx context.Context, q *queries.Queries, qIntermediateSession queries.IntermediateSession, qUser queries.User) error {
	userHasPasskey, err := q.GetUserHasActivePasskey(ctx, qUser.ID)
	if err != nil {
		return fmt.Errorf("get user has passkey: %w", err)
	}

	if userHasPasskey {
		return fmt.Errorf("user already has a passkey")
	}

	if _, err := q.CreatePasskey(ctx, queries.CreatePasskeyParams{
		ID:           uuid.New(),
		UserID:       qUser.ID,
		CredentialID: qIntermediateSession.PasskeyCredentialID,
		PublicKey:    qIntermediateSession.PasskeyPublicKey,
		Aaguid:       *qIntermediateSession.PasskeyAaguid,
		RpID:         *qIntermediateSession.PasskeyRpID,
	}); err != nil {
		return fmt.Errorf("create passkey: %w", err)
	}

	return nil
}

func (s *Store) copyRegisteredAuthenticatorAppSettings(ctx context.Context, q *queries.Queries, qIntermediateSession queries.IntermediateSession, qUser queries.User) error {
	if qUser.AuthenticatorAppSecretCiphertext != nil || qUser.AuthenticatorAppRecoveryCodeBcrypts != nil {
		return fmt.Errorf("user already has authenticator app registered")
	}

	if _, err := q.UpdateUserAuthenticatorApp(ctx, queries.UpdateUserAuthenticatorAppParams{
		AuthenticatorAppSecretCiphertext:    qIntermediateSession.AuthenticatorAppSecretCiphertext,
		AuthenticatorAppRecoveryCodeBcrypts: qIntermediateSession.AuthenticatorAppRecoveryCodeBcrypts,
		ID:                                  qUser.ID,
	}); err != nil {
		return fmt.Errorf("update user authenticator app: %w", err)
	}

	return nil
}
