package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

const sessionDuration = time.Hour * 24 * 7
const relayedSessionTokenDuration = time.Minute

func (s *Store) ExchangeIntermediateSessionForSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForSessionResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
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

	qProjectUISettings, err := q.GetProjectUISettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project ui settings: %w", err)
	}

	slog.InfoContext(ctx, "exchange_intermediate_session_for_session",
		"organization_id", idformat.Organization.Format(qOrg.ID),
		"intermediate_session_id", idformat.IntermediateSession.Format(qIntermediateSession.ID),
		"email", qIntermediateSession.Email)

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

	var (
		newUser        = qUser == nil
		detailsUpdated = newUser
	)

	// if no matching user, create a new one
	if qUser == nil {
		if !qProjectUISettings.SelfServeCreateUsers {
			return nil, apierror.NewPermissionDeniedError("self-serve user creation disabled", nil)
		}

		slog.InfoContext(ctx, "create_user")
		qNewUser, err := q.CreateUser(ctx, queries.CreateUserParams{
			ID:                uuid.New(),
			OrganizationID:    qOrg.ID,
			Email:             *qIntermediateSession.Email,
			DisplayName:       qIntermediateSession.UserDisplayName,
			ProfilePictureUrl: qIntermediateSession.ProfilePictureUrl,
			GoogleUserID:      qIntermediateSession.GoogleUserID,
			MicrosoftUserID:   qIntermediateSession.MicrosoftUserID,
			GithubUserID:      qIntermediateSession.GithubUserID,
			PasswordBcrypt:    qIntermediateSession.NewUserPasswordBcrypt,
		})
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}

		auditUser, err := s.auditlogStore.GetUser(ctx, tx, qNewUser.ID)
		if err != nil {
			return nil, fmt.Errorf("get audit user: %w", err)
		}

		if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
			EventName: "tesseral.users.create",
			EventDetails: &auditlogv1.CreateUser{
				User: auditUser,
			},
			OrganizationID: &qOrg.ID,
			ResourceType:   queries.AuditLogEventResourceTypeUser,
			ResourceID:     &qNewUser.ID,
		}); err != nil {
			return nil, fmt.Errorf("log audit event: %w", err)
		}

		qUser = &qNewUser
	} else {
		detailsUpdated =
			(qIntermediateSession.GithubUserID != nil && *qIntermediateSession.GithubUserID != derefOrEmpty(qUser.GithubUserID)) ||
				(qIntermediateSession.GoogleUserID != nil && *qIntermediateSession.GoogleUserID != derefOrEmpty(qUser.GoogleUserID)) ||
				(qIntermediateSession.MicrosoftUserID != nil && *qIntermediateSession.MicrosoftUserID != derefOrEmpty(qUser.MicrosoftUserID)) ||
				(qIntermediateSession.UserDisplayName != nil && *qIntermediateSession.UserDisplayName != derefOrEmpty(qUser.DisplayName)) ||
				(qIntermediateSession.ProfilePictureUrl != nil && *qIntermediateSession.ProfilePictureUrl != derefOrEmpty(qUser.ProfilePictureUrl)) ||
				qIntermediateSession.NewUserPasswordBcrypt != nil

		if detailsUpdated {
			slog.InfoContext(ctx, "update_user")

			auditPreviousUser, err := s.auditlogStore.GetUser(ctx, tx, qUser.ID)
			if err != nil {
				return nil, fmt.Errorf("get audit user: %w", err)
			}

			qUpdatedUser, err := q.UpdateUserDetails(ctx, queries.UpdateUserDetailsParams{
				ID:                qUser.ID,
				GithubUserID:      qIntermediateSession.GithubUserID,
				GoogleUserID:      qIntermediateSession.GoogleUserID,
				MicrosoftUserID:   qIntermediateSession.MicrosoftUserID,
				DisplayName:       qIntermediateSession.UserDisplayName,
				ProfilePictureUrl: qIntermediateSession.ProfilePictureUrl,
				PasswordBcrypt:    qIntermediateSession.NewUserPasswordBcrypt,
			})
			if err != nil {
				return nil, fmt.Errorf("update user: %w", err)
			}

			auditUser, err := s.auditlogStore.GetUser(ctx, tx, qUser.ID)
			if err != nil {
				return nil, fmt.Errorf("get audit user: %w", err)
			}

			if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
				EventName: "tesseral.users.update",
				EventDetails: &auditlogv1.UpdateUser{
					User:         auditUser,
					PreviousUser: auditPreviousUser,
				},
				OrganizationID: &qOrg.ID,
				ResourceType:   queries.AuditLogEventResourceTypeUser,
				ResourceID:     &qUser.ID,
			}); err != nil {
				return nil, fmt.Errorf("log audit event: %w", err)
			}

			qUser = &qUpdatedUser
		}
	}

	// if a passkey is registered on the intermediate session, copy it onto the
	// user
	if qIntermediateSession.PasskeyCredentialID != nil {
		slog.InfoContext(ctx, "register_passkey")
		detailsUpdated = true
		if err := s.copyRegisteredPasskeySettings(ctx, q, qIntermediateSession, *qUser); err != nil {
			return nil, fmt.Errorf("copy registered passkey settings: %w", err)
		}
	}

	// if an authenticator app is registered on the intermediate session, copy
	// it onto the user
	if qIntermediateSession.AuthenticatorAppSecretCiphertext != nil {
		slog.InfoContext(ctx, "register_authenticator_app")
		detailsUpdated = true
		if err := s.copyRegisteredAuthenticatorAppSettings(ctx, q, qIntermediateSession, *qUser); err != nil {
			return nil, fmt.Errorf("copy registered authenticator app settings: %w", err)
		}
	}

	expireTime := time.Now().Add(sessionDuration)

	// Create a new session for the user
	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])
	qSession, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expireTime,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qUser.ID,
		PrimaryAuthFactor:  *qIntermediateSession.PrimaryAuthFactor,
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
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
		slog.InfoContext(ctx, "upgrade_user_invite", "is_owner", qUserInvite.IsOwner)
		if qUserInvite.IsOwner {
			detailsUpdated = true
			if _, err := q.UpdateUserIsOwner(ctx, queries.UpdateUserIsOwnerParams{
				ID:      qUser.ID,
				IsOwner: true,
			}); err != nil {
				return nil, fmt.Errorf("update user is owner: %w", err)
			}
		}
	}

	// if the intermediate session has a relayed session state, create a relayed
	// session
	var relayedSessionToken string
	if qIntermediateSession.RelayedSessionState != nil {
		slog.InfoContext(ctx, "create_relayed_session")
		relayedSessionTokenUUID := uuid.New()
		relayedSessionTokenSHA256 := sha256.Sum256(relayedSessionTokenUUID[:])
		relayedSessionTokenExpireTime := time.Now().Add(relayedSessionTokenDuration)

		if _, err := q.CreateRelayedSession(ctx, queries.CreateRelayedSessionParams{
			SessionID:                     qSession.ID,
			RelayedSessionTokenExpireTime: &relayedSessionTokenExpireTime,
			RelayedSessionTokenSha256:     relayedSessionTokenSHA256[:],
			RelayedRefreshTokenSha256:     nil, // assigned in ExchangeRelayedSessionTokenForSession
			State:                         qIntermediateSession.RelayedSessionState,
		}); err != nil {
			return nil, fmt.Errorf("create relayed session: %w", err)
		}

		relayedSessionToken = idformat.RelayedSessionToken.Format(relayedSessionTokenUUID)
	}

	auditSession, err := s.auditlogStore.GetSession(ctx, tx, qSession.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit session: %w", err)
	}

	var samlConnectionID *string
	if qIntermediateSession.VerifiedSamlConnectionID != nil {
		samlConnectionID = refOrNil(idformat.SAMLConnection.Format(*qIntermediateSession.VerifiedSamlConnectionID))
	}

	var oidcConnectionID *string
	if qIntermediateSession.VerifiedOidcConnectionID != nil {
		oidcConnectionID = refOrNil(idformat.OIDCConnection.Format(*qIntermediateSession.VerifiedOidcConnectionID))
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.sessions.create",
		EventDetails: &auditlogv1.CreateSession{
			Session:          auditSession,
			SamlConnectionId: samlConnectionID,
			OidcConnectionId: oidcConnectionID,
		},
		OrganizationID: &qOrg.ID,
		ResourceType:   queries.AuditLogEventResourceTypeSession,
		ResourceID:     &qSession.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	if detailsUpdated {
		// Send sync user event
		if err := s.sendSyncUserEvent(ctx, *qUser); err != nil {
			return nil, fmt.Errorf("send sync user event: %w", err)
		}
	}

	return &intermediatev1.ExchangeIntermediateSessionForSessionResponse{
		AccessToken:                           "", // populated in service
		RefreshToken:                          idformat.SessionRefreshToken.Format(refreshToken),
		NewUser:                               newUser,
		RelayedSessionToken:                   relayedSessionToken,
		RedirectUri:                           derefOrEmpty(qIntermediateSession.RedirectUri),
		ReturnRelayedSessionTokenAsQueryParam: qIntermediateSession.ReturnRelayedSessionTokenAsQueryParam,
	}, nil
}

func (s *Store) sendSyncUserEvent(ctx context.Context, qUser queries.User) error {
	qProjectWebhookSettings, err := s.q.GetProjectWebhookSettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		// We want to ignore this error if the project does not have webhook settings
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("get project by id: %w", err)
	}

	message, err := s.svixClient.Message.Create(ctx, qProjectWebhookSettings.AppID, models.MessageIn{
		EventType: "sync.user",
		Payload: map[string]interface{}{
			"type":   "sync.user",
			"userId": idformat.User.Format(qUser.ID),
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	slog.InfoContext(ctx, "svix_message_created", "message_id", message.Id, "event_type", message.EventType, "user_id", idformat.User.Format(qUser.ID))

	return nil
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

	if qIntermediateSession.PrimaryAuthFactor == nil {
		return apierror.NewFailedPreconditionError("primary auth factor not set", nil)
	}

	if !emailVerified {
		return apierror.NewFailedPreconditionError("email not verified", nil)
	}

	isEnterpriseLogin :=
		*qIntermediateSession.PrimaryAuthFactor == queries.PrimaryAuthFactorSaml ||
			*qIntermediateSession.PrimaryAuthFactor == queries.PrimaryAuthFactorOidc

	if !isEnterpriseLogin {
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
	}

	switch *qIntermediateSession.PrimaryAuthFactor {
	case queries.PrimaryAuthFactorEmail:
		if qOrg.LogInWithEmail {
			return nil
		}
	case queries.PrimaryAuthFactorPassword:
		if qOrg.LogInWithPassword {
			return nil
		}
	case queries.PrimaryAuthFactorGoogle:
		if qIntermediateSession.GoogleUserID == nil {
			panic(fmt.Errorf("intermediate session missing google user id: %v", qIntermediateSession.ID))
		}

		if qOrg.LogInWithGoogle {
			return nil
		}
	case queries.PrimaryAuthFactorMicrosoft:
		if qIntermediateSession.MicrosoftUserID == nil {
			panic(fmt.Errorf("intermediate session missing microsoft user id: %v", qIntermediateSession.ID))
		}

		if qOrg.LogInWithMicrosoft {
			return nil
		}
	case queries.PrimaryAuthFactorGithub:
		if qIntermediateSession.GithubUserID == nil {
			panic(fmt.Errorf("intermediate session missing github user id: %v", qIntermediateSession.ID))
		}

		if qOrg.LogInWithGithub {
			return nil
		}
	case queries.PrimaryAuthFactorSaml:
		if qIntermediateSession.VerifiedSamlConnectionID == nil {
			panic(fmt.Errorf("intermediate session missing verified saml connection id: %v", qIntermediateSession.ID))
		}
		if qOrg.LogInWithSaml {
			return nil
		}
	case queries.PrimaryAuthFactorOidc:
		if qIntermediateSession.VerifiedOidcConnectionID == nil {
			panic(fmt.Errorf("intermediate session missing verified oidc connection id: %v", qIntermediateSession.ID))
		}
		if qOrg.LogInWithOidc {
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

	qUser, err = s.matchGithubUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match github user: %w", err)
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

func (s *Store) matchGithubUser(ctx context.Context, q *queries.Queries, qOrg queries.Organization, qIntermediateSession queries.IntermediateSession) (*queries.User, error) {
	if qIntermediateSession.GithubUserID == nil {
		return nil, nil
	}

	qUser, err := q.GetOrganizationUserByGithubUserID(ctx, queries.GetOrganizationUserByGithubUserIDParams{
		OrganizationID: qOrg.ID,
		GithubUserID:   qIntermediateSession.GithubUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get organization user by github user id: %w", err)
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
	if qUser.AuthenticatorAppSecretCiphertext != nil || qUser.AuthenticatorAppRecoveryCodeSha256s != nil {
		return fmt.Errorf("user already has authenticator app registered")
	}

	if _, err := q.UpdateUserAuthenticatorApp(ctx, queries.UpdateUserAuthenticatorAppParams{
		AuthenticatorAppSecretCiphertext:    qIntermediateSession.AuthenticatorAppSecretCiphertext,
		AuthenticatorAppRecoveryCodeSha256s: qIntermediateSession.AuthenticatorAppRecoveryCodeSha256s,
		ID:                                  qUser.ID,
	}); err != nil {
		return fmt.Errorf("update user authenticator app: %w", err)
	}

	return nil
}
