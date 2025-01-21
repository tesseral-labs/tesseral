package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	openauthecdsa "github.com/openauth/openauth/internal/crypto/ecdsa"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/sessions"
	"github.com/openauth/openauth/internal/store/idformat"
)

var errInvalidIntermediateSessionState = fmt.Errorf("invalid intermediate session state")

func (s *Store) ExchangeIntermediateSessionForNewOrganizationSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	projectID := authn.ProjectID(ctx)

	// Create a new organization
	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                   uuid.New(),
		ProjectID:            projectID,
		DisplayName:          req.DisplayName,
		GoogleHostedDomain:   refOrNil(intermediateSession.GoogleHostedDomain),
		MicrosoftTenantID:    refOrNil(intermediateSession.MicrosoftTenantId),
		OverrideLogInMethods: false,
	})
	if err != nil {
		return nil, err
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

	project, err := q.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * time.Hour * 24) // 7 days

	// Create a new session for the user
	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])

	qSession, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expiresAt,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qUser.ID,
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
		ID:          idformat.Organization.Format(qOrganization.ID),
		DisplayName: qOrganization.DisplayName,
	}, &sessions.Project{
		ID: idformat.Project.Format(project.ID),
	}, &sessions.Session{
		ID:         idformat.Session.Format(qSession.ID),
		UserID:     idformat.User.Format(qUser.ID),
		CreateTime: *qSession.CreateTime,
		ExpireTime: *qSession.ExpireTime,
		Revoked:    qSession.Revoked,
	}, &sessions.User{
		ID:              idformat.User.Format(qUser.ID),
		CreateTime:      *qUser.CreateTime,
		Email:           qUser.Email,
		GoogleUserID:    derefOrEmpty(qUser.GoogleUserID),
		MicrosoftUserID: derefOrEmpty(qUser.MicrosoftUserID),
		UpdateTime:      derefOrEmpty(qUser.UpdateTime),
	}, *sessionSigningKeyID, privateKey)
	if err != nil {
		return nil, err
	}

	return &intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func (s *Store) ExchangeIntermediateSessionForSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	projectID := authn.ProjectID(ctx)

	organizationID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qOrganization, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ID:        organizationID,
		ProjectID: projectID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get project organization by id: %w", err))
		}

		return nil, fmt.Errorf("get project organization by id: %w", err)
	}

	// Use the intermediate session state to determine the user to sign in
	// The hierarchy of user identifiers is:
	// 1. Microsoft user ID
	// 2. Google user ID
	// 3. Email
	var qUser queries.User
	if intermediateSession.MicrosoftUserId != "" {
		qUser, err = q.GetOrganizationUserByMicrosoftUserID(ctx, queries.GetOrganizationUserByMicrosoftUserIDParams{
			OrganizationID:  organizationID,
			MicrosoftUserID: &intermediateSession.MicrosoftUserId,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get organization user by microsoft user id: %w", err))
			}

			return nil, fmt.Errorf("get organization user by microsoft user id: %w", err)
		}
	} else if intermediateSession.GoogleUserId != "" {
		qUser, err = q.GetOrganizationUserByGoogleUserID(ctx, queries.GetOrganizationUserByGoogleUserIDParams{
			OrganizationID: organizationID,
			GoogleUserID:   &intermediateSession.GoogleUserId,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get organization user by google user id: %w", err))
			}

			return nil, fmt.Errorf("get organization user by google user id: %w", err)
		}
	} else if intermediateSession.Email != "" {
		qUser, err = q.GetOrganizationUserByEmail(ctx, queries.GetOrganizationUserByEmailParams{
			OrganizationID: organizationID,
			Email:          intermediateSession.Email,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get organization user by email: %w", err))
			}

			return nil, fmt.Errorf("get organization user by email: %w", err)
		}

		// Ensure that the intermediate session is in an authorized state
		if !intermediateSession.PasswordVerified || intermediateSession.OrganizationId != req.OrganizationId {
			return nil, fmt.Errorf("verify intermediate session state: %w", errInvalidIntermediateSessionState)
		}
	}

	expiresAt := time.Now().Add(7 * time.Hour * 24) // 7 days

	// Create a new session for the user
	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])

	qSession, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expiresAt,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qUser.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	sessionSigningKeyID, privateKey, err := s.getSessionSigningKey(ctx, q, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("session signing key not found", fmt.Errorf("get current session key by project id: %w", err))
		}

		return nil, fmt.Errorf("get session signing key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	accessToken, err := sessions.GetAccessToken(ctx, &sessions.Organization{
		ID:          idformat.Organization.Format(qOrganization.ID),
		DisplayName: qOrganization.DisplayName,
	}, &sessions.Project{
		ID: idformat.Project.Format(qProject.ID),
	}, &sessions.Session{
		ID:         idformat.Session.Format(qSession.ID),
		UserID:     idformat.User.Format(qUser.ID),
		CreateTime: *qSession.CreateTime,
		ExpireTime: *qSession.ExpireTime,
		Revoked:    qSession.Revoked,
	}, &sessions.User{
		ID:              idformat.User.Format(qUser.ID),
		CreateTime:      *qUser.CreateTime,
		Email:           qUser.Email,
		GoogleUserID:    derefOrEmpty(qUser.GoogleUserID),
		MicrosoftUserID: derefOrEmpty(qUser.MicrosoftUserID),
		UpdateTime:      derefOrEmpty(qUser.UpdateTime),
	}, *sessionSigningKeyID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	return &intermediatev1.ExchangeIntermediateSessionForSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func (s *Store) getSessionSigningKey(ctx context.Context, q *queries.Queries, projectID uuid.UUID) (*uuid.UUID, *ecdsa.PrivateKey, error) {
	sessionSigningKey, err := q.GetCurrentSessionKeyByProjectID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, apierror.NewFailedPreconditionError("session signing key not found", fmt.Errorf("get current session key by project id: %w", err))
		}

		return nil, nil, fmt.Errorf("get current session key by project id: %w", err)
	}

	decryptResult, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      sessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt: %w", err)
	}

	privateKey, err := openauthecdsa.PrivateKeyFromBytes(decryptResult.Plaintext)
	if err != nil {
		return nil, nil, fmt.Errorf("private key from bytes: %w", err)
	}

	return &sessionSigningKey.ID, privateKey, nil
}
