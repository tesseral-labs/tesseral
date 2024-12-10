package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	openauthecdsa "github.com/openauth/openauth/internal/crypto/ecdsa"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/sessions"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) ExchangeIntermediateSessionForNewOrganizationSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	projectID := projectid.ProjectID(ctx)

	// Create a new organization
	organization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                 uuid.New(),
		ProjectID:          projectID,
		DisplayName:        req.DisplayName,
		GoogleHostedDomain: &intermediateSession.GoogleHostedDomain,
		MicrosoftTenantID:  &intermediateSession.MicrosoftTenantId,
	})
	if err != nil {
		return nil, err
	}

	// Create a new user for that organization
	user, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:              uuid.New(),
		OrganizationID:  organization.ID,
		Email:           intermediateSession.Email,
		GoogleUserID:    &intermediateSession.GoogleUserId,
		MicrosoftUserID: &intermediateSession.MicrosoftUserId,
	})
	if err != nil {
		return nil, err
	}

	slog.Info("ExchangeIntermediateSessionForNewOrganizationSession", "organization", organization, "user", user)

	project, err := q.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * time.Hour * 24) // 7 days

	// Create a new session for the user
	token := uuid.New()
	refreshToken := idformat.SessionRefreshToken.Format(token)
	refreshTokenSha256 := sha256.Sum256([]byte(refreshToken))

	session, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.New(),
		ExpireTime:         &expiresAt,
		RefreshTokenSha256: refreshTokenSha256[:],
		UserID:             user.ID,
	})
	if err != nil {
		return nil, err
	}

	sessionSigningKeyID, privateKey, err := s.getSessionSigningKey(ctx, q, projectID)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	accessToken, err := sessions.GetAccessToken(ctx, &sessions.Organization{
		ID:          idformat.Organization.Format(organization.ID),
		DisplayName: organization.DisplayName,
	}, &sessions.Project{
		ID: idformat.Project.Format(project.ID),
	}, &sessions.Session{
		ID:         idformat.Session.Format(session.ID),
		UserID:     idformat.User.Format(user.ID),
		CreateTime: *session.CreateTime,
		ExpireTime: *session.ExpireTime,
		Revoked:    session.Revoked,
	}, &sessions.User{
		ID:              idformat.User.Format(user.ID),
		CreateTime:      user.CreateTime.Time,
		Email:           user.Email,
		GoogleUserID:    *user.GoogleUserID,
		MicrosoftUserID: *user.MicrosoftUserID,
		UpdateTime:      user.UpdateTime.Time,
	}, *sessionSigningKeyID, privateKey)
	if err != nil {
		return nil, err
	}

	slog.Info("ExchangeIntermediateSessionForNewOrganizationSession", "accessToken", accessToken)

	return &intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (s *Store) ExchangeIntermediateSessionForSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	projectID := projectid.ProjectID(ctx)

	organizationID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	organization, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ID:        organizationID,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}

	slog.Info("ExchangeIntermediateSessionForSession", "organization", organization)

	// Use the intermediate session state to determine the user to sign in
	// The hierarchy of user identifiers is:
	// 1. Microsoft user ID
	// 2. Google user ID
	// 3. Email
	var user queries.User
	if intermediateSession.MicrosoftUserId != "" {
		user, err = q.GetOrganizationUserByMicrosoftUserID(ctx, queries.GetOrganizationUserByMicrosoftUserIDParams{
			OrganizationID:  organizationID,
			MicrosoftUserID: &intermediateSession.MicrosoftUserId,
		})
		if err != nil {
			return nil, err
		}
	} else if intermediateSession.GoogleUserId != "" {
		user, err = q.GetOrganizationUserByGoogleUserID(ctx, queries.GetOrganizationUserByGoogleUserIDParams{
			OrganizationID: organizationID,
			GoogleUserID:   &intermediateSession.GoogleUserId,
		})
		if err != nil {
			return nil, err
		}
	} else if intermediateSession.Email != "" {
		user, err = q.GetOrganizationUserByEmail(ctx, queries.GetOrganizationUserByEmailParams{
			OrganizationID: organizationID,
			Email:          intermediateSession.Email,
		})
		if err != nil {
			return nil, err
		}
	}

	slog.Info("ExchangeIntermediateSessionForSession", "user", user)

	expiresAt := time.Now().Add(7 * time.Hour * 24) // 7 days

	// Create a new session for the user
	token := uuid.New()
	refreshToken := idformat.SessionRefreshToken.Format(token)
	refreshTokenSha256 := sha256.Sum256([]byte(refreshToken))

	session, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.New(),
		ExpireTime:         &expiresAt,
		RefreshTokenSha256: refreshTokenSha256[:],
		UserID:             user.ID,
	})
	if err != nil {
		return nil, err
	}

	project, err := q.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	sessionSigningKeyID, privateKey, err := s.getSessionSigningKey(ctx, q, projectID)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	accessToken, err := sessions.GetAccessToken(ctx, &sessions.Organization{
		ID:          idformat.Organization.Format(organization.ID),
		DisplayName: organization.DisplayName,
	}, &sessions.Project{
		ID: idformat.Project.Format(project.ID),
	}, &sessions.Session{
		ID:         idformat.Session.Format(session.ID),
		UserID:     idformat.User.Format(user.ID),
		CreateTime: *session.CreateTime,
		ExpireTime: *session.ExpireTime,
		Revoked:    session.Revoked,
	}, &sessions.User{
		ID:              idformat.User.Format(user.ID),
		CreateTime:      user.CreateTime.Time,
		Email:           user.Email,
		GoogleUserID:    *user.GoogleUserID,
		MicrosoftUserID: *user.MicrosoftUserID,
		UpdateTime:      user.UpdateTime.Time,
	}, *sessionSigningKeyID, privateKey)
	if err != nil {
		return nil, err
	}

	slog.Info("ExchangeIntermediateSessionForSession", "accessToken", accessToken)

	return &intermediatev1.ExchangeIntermediateSessionForSessionResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (s *Store) getSessionSigningKey(ctx context.Context, q *queries.Queries, projectID uuid.UUID) (*uuid.UUID, *ecdsa.PrivateKey, error) {
	sessionSigningKey, err := q.GetCurrentSessionKeyByProjectID(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}

	decryptResult, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      sessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
	})
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := openauthecdsa.PrivateKeyFromBytes(decryptResult.Value)
	if err != nil {
		return nil, nil, err
	}

	return &sessionSigningKey.ID, privateKey, nil
}
