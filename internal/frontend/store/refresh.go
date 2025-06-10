package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) LogRefreshEvent(ctx context.Context, refreshToken string) error {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer rollback()

	var qDetails struct {
		SessionID               uuid.UUID
		UserID                  uuid.UUID
		OrganizationID          uuid.UUID
		UserIsOwner             bool
		UserEmail               string
		UserDisplayName         *string
		UserProfilePictureUrl   *string
		OrganizationDisplayName string
		ImpersonatorUserID      *uuid.UUID
		ProjectID               uuid.UUID
	}

	switch {
	case strings.HasPrefix(refreshToken, "tesseral_secret_session_refresh_token_"):
		slog.InfoContext(ctx, "refresh_session_token")

		refreshTokenUUID, err := idformat.SessionRefreshToken.Parse(refreshToken)
		if err != nil {
			return fmt.Errorf("parse refresh token: %w", err)
		}

		refreshTokenSHA := sha256.Sum256(refreshTokenUUID[:])
		qSessionDetails, err := q.GetSessionDetailsByRefreshTokenSHA256(ctx, refreshTokenSHA[:])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return apierror.NewUnauthenticatedError("invalid refresh token", fmt.Errorf("invalid refresh token"))
			}

			return fmt.Errorf("get session details by refresh token sha256: %w", err)
		}

		qDetails.SessionID = qSessionDetails.SessionID
		qDetails.UserID = qSessionDetails.UserID
		qDetails.OrganizationID = qSessionDetails.OrganizationID
		qDetails.ProjectID = qSessionDetails.ProjectID
	case strings.HasPrefix(refreshToken, "tesseral_secret_relayed_session_refresh_token_"):
		slog.InfoContext(ctx, "refresh_relayed_session_token")

		relayedRefreshTokenUUID, err := idformat.RelayedSessionRefreshToken.Parse(refreshToken)
		if err != nil {
			return fmt.Errorf("parse refresh token: %w", err)
		}

		relayedRefreshTokenSHA := sha256.Sum256(relayedRefreshTokenUUID[:])
		qSessionDetails, err := q.GetSessionDetailsByRelayedSessionRefreshTokenSHA256(ctx, relayedRefreshTokenSHA[:])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return apierror.NewUnauthenticatedError("invalid refresh token", fmt.Errorf("invalid refresh token"))
			}

			return fmt.Errorf("get session details by refresh token sha256: %w", err)
		}

		qDetails.SessionID = qSessionDetails.SessionID
		qDetails.UserID = qSessionDetails.UserID
		qDetails.OrganizationID = qSessionDetails.OrganizationID
		qDetails.ProjectID = qSessionDetails.ProjectID
	}

	qSession, err := q.GetSessionByID(ctx, qDetails.SessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.NewUnauthenticatedError("invalid session", fmt.Errorf("invalid session"))
		}
		return fmt.Errorf("get session: %w", err)
	}

	var qImpersonatingUserEmail *string
	if qDetails.ImpersonatorUserID != nil {
		qImpersonatingUser, err := q.GetUserByID(ctx, *qDetails.ImpersonatorUserID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("get impersonator user: %w", err)
			}
		}

		qImpersonatingUserEmail = &qImpersonatingUser.Email
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.sessions.refresh",
		EventDetails: map[string]any{
			"session": parseSessionEventDetails(qSession, qImpersonatingUserEmail),
		},
		OrganizationID: &qDetails.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeSession,
		ResourceID:     &qDetails.SessionID,
	}); err != nil {
		return fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
