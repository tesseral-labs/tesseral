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

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		return nil, err
	}

	if err := enforceOrganizationLoginEnabled(qOrg); err != nil {
		return nil, fmt.Errorf("enforce organization login enabled: %w", err)
	}

	if err := s.validateAuthRequirementsSatisfied(ctx, q, qIntermediateSession.ID, qOrg.ID); err != nil {
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
		})
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}

		qUser = &qNewUser
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

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.ExchangeIntermediateSessionForSessionResponse{
		AccessToken:  "", // populated in service
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func (s *Store) ExchangeIntermediateSessionForNewOrganizationSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("get project by id: %w", fmt.Errorf("project not found: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	// Create a new organization
	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                   uuid.New(),
		ProjectID:            authn.ProjectID(ctx),
		DisplayName:          req.DisplayName,
		OverrideLogInMethods: false,
	})
	if err != nil {
		return nil, err
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

	// Create a new user for that organization
	qUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:              uuid.New(),
		OrganizationID:  qOrganization.ID,
		Email:           intermediateSession.Email,
		GoogleUserID:    refOrNil(intermediateSession.GoogleUserId),
		MicrosoftUserID: refOrNil(intermediateSession.MicrosoftUserId),
		IsOwner:         true,
	})
	if err != nil {
		return nil, err
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
	intermediateSessionUUID, err := idformat.IntermediateSession.Parse(intermediateSession.Id)
	if err != nil {
		panic(fmt.Errorf("parse intermediate session id: %w", err))
	}

	if _, err := q.RevokeIntermediateSession(ctx, intermediateSessionUUID); err != nil {
		return nil, fmt.Errorf("revoke intermediate session: %w", err)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse{
		AccessToken:  "", // populated in service
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func (s *Store) validateAuthRequirementsSatisfied(ctx context.Context, q *queries.Queries, intermediateSessionID, organizationID uuid.UUID) error {
	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, intermediateSessionID)
	if err != nil {
		return fmt.Errorf("get intermediate session by id: %w", err)
	}

	emailVerified, err := s.getIntermediateSessionEmailVerified(ctx, q, qIntermediateSession.ID)
	if err != nil {
		return fmt.Errorf("get intermediate session verified: %w", err)
	}

	if !emailVerified {
		return apierror.NewFailedPreconditionError("email not verified", nil)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        organizationID,
	})
	if err != nil {
		return fmt.Errorf("get organization by id: %w", err)
	}

	googleEnabled := qProject.LogInWithGoogleEnabled
	if derefOrEmpty(qOrg.DisableLogInWithGoogle) {
		googleEnabled = false
	}

	microsoftEnabled := qProject.LogInWithMicrosoftEnabled
	if derefOrEmpty(qOrg.DisableLogInWithMicrosoft) {
		microsoftEnabled = false
	}

	passwordEnabled := qProject.LogInWithPasswordEnabled
	if derefOrEmpty(qOrg.DisableLogInWithPassword) {
		passwordEnabled = false
	}

	if googleEnabled && qIntermediateSession.GoogleUserID != nil {
		return nil
	}
	if microsoftEnabled && qIntermediateSession.MicrosoftUserID != nil {
		return nil
	}
	if passwordEnabled && qIntermediateSession.PasswordVerified != nil && *qIntermediateSession.PasswordVerified { // TODO password_verified should be non-null
		return nil
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
