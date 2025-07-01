package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Store) RedeemUserImpersonationToken(ctx context.Context, req *intermediatev1.RedeemUserImpersonationTokenRequest) (*intermediatev1.RedeemUserImpersonationTokenResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	secretTokenUUID, err := idformat.UserImpersonationSecretToken.Parse(req.SecretUserImpersonationToken)
	if err != nil {
		return nil, fmt.Errorf("parse user impersonation secret token: %w", err)
	}

	secretTokenSHA256 := sha256.Sum256(secretTokenUUID[:])
	qUserImpersonationToken, err := q.GetUserImpersonationTokenBySecretTokenSHA256(ctx, secretTokenSHA256[:])
	if err != nil {
		return nil, fmt.Errorf("get user impersonation token by secret token sha256: %w", err)
	}

	expireTime := time.Now().Add(sessionDuration)

	// Create a new session for the user
	slog.InfoContext(ctx, "impersonate_user",
		"impersonator_id", idformat.User.Format(qUserImpersonationToken.ImpersonatorID),
		"impersonated_id", idformat.User.Format(qUserImpersonationToken.ImpersonatedID))

	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])

	qSession, err := q.CreateImpersonatedSession(ctx, queries.CreateImpersonatedSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expireTime,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qUserImpersonationToken.ImpersonatedID,
		ImpersonatorUserID: &qUserImpersonationToken.ImpersonatorID,
	})
	if err != nil {
		return nil, fmt.Errorf("create impersonated session: %w", err)
	}

	if _, err := q.RevokeUserImpersonationToken(ctx, qUserImpersonationToken.ID); err != nil {
		return nil, fmt.Errorf("revoke user impersonation token: %w", err)
	}

	// Get impersonated user for logging purposes
	qImpersonatedUser, err := q.GetUserByID(ctx, qSession.UserID)
	if err != nil {
		return nil, fmt.Errorf("get impersonated user by id: %w", err)
	}

	auditSession, err := s.auditlogStore.GetSession(ctx, tx, qSession.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit session: %w", err)
	}

	if err := logImpersonateAuditEvent(ctx, q, logAuditEventParams{
		EventName:      "tesseral.sessions.create",
		EventDetails:   &auditlogv1.CreateSession{Session: auditSession},
		OrganizationID: &qImpersonatedUser.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeSession,
		ResourceID:     &qSession.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RedeemUserImpersonationTokenResponse{
		AccessToken:  "", // populated in service
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func logImpersonateAuditEvent(ctx context.Context, q *queries.Queries, data logAuditEventParams) error {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now()
	eventID := uuidv7.NewWithTime(eventTime)

	eventDetailsBytes, err := protojson.Marshal(data.EventDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal event details: %w", err)
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: data.OrganizationID,
		ResourceType:   &data.ResourceType,
		ResourceID:     data.ResourceID,
		EventName:      data.EventName,
		EventTime:      &eventTime,
		EventDetails:   eventDetailsBytes,
	}

	_, err = q.CreateAuditLogEvent(ctx, qEventParams)
	return err
}
