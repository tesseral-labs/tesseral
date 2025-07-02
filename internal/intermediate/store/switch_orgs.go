package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) ExchangeSessionForIntermediateSession(ctx context.Context, req *intermediatev1.ExchangeSessionForIntermediateSessionRequest) (*intermediatev1.ExchangeSessionForIntermediateSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

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

	refreshTokenUUID, err := idformat.SessionRefreshToken.Parse(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}

	refreshTokenSHA := sha256.Sum256(refreshTokenUUID[:])
	qDetails, err := s.q.GetSessionDetailsByRefreshTokenSHA256(ctx, queries.GetSessionDetailsByRefreshTokenSHA256Params{
		ProjectID:          authn.ProjectID(ctx),
		RefreshTokenSha256: refreshTokenSHA[:],
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewUnauthenticatedError("invalid refresh token", fmt.Errorf("invalid refresh token"))
		}

		return nil, fmt.Errorf("get session details by refresh token sha256: %w", err)
	}

	var primaryAuthFactor *queries.PrimaryAuthFactor
	var googleUserID, microsoftUserID *string
	switch qDetails.PrimaryAuthFactor {
	case queries.PrimaryAuthFactorEmail:
		primaryAuthFactor = &qDetails.PrimaryAuthFactor
	case queries.PrimaryAuthFactorGoogle:
		primaryAuthFactor = &qDetails.PrimaryAuthFactor
		googleUserID = qDetails.GoogleUserID
	case queries.PrimaryAuthFactorMicrosoft:
		primaryAuthFactor = &qDetails.PrimaryAuthFactor
		microsoftUserID = qDetails.MicrosoftUserID
	}

	expireTime := time.Now().Add(intermediateSessionDuration)

	secretToken := uuid.New()
	secretTokenSHA256 := sha256.Sum256(secretToken[:])
	if _, err := q.CreateIntermediateSession(ctx, queries.CreateIntermediateSessionParams{
		ID:                uuid.Must(uuid.NewV7()),
		ProjectID:         authn.ProjectID(ctx),
		ExpireTime:        &expireTime,
		Email:             &qDetails.Email,
		GoogleUserID:      googleUserID,
		MicrosoftUserID:   microsoftUserID,
		SecretTokenSha256: secretTokenSHA256[:],
		PrimaryAuthFactor: primaryAuthFactor,

		// If the input session was from a "Log in with Email", then carry over
		// that session's email verification to the new intermediate session.
		//
		// For other login methods, the existing oauth_verified_emails entry
		// will have this effect already.
		EmailVerificationChallengeCompleted: qDetails.PrimaryAuthFactor == queries.PrimaryAuthFactorEmail,
	}); err != nil {
		return nil, fmt.Errorf("create intermediate session: %w", err)
	}

	// revoke input session
	if err := q.InvalidateSession(ctx, qDetails.SessionID); err != nil {
		return nil, fmt.Errorf("invalidate session: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.ExchangeSessionForIntermediateSessionResponse{
		IntermediateSessionSecretToken: idformat.IntermediateSessionSecretToken.Format(secretToken),
	}, nil
}
