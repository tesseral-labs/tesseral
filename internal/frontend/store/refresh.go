package store

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Store) CreateRefreshAuditLogEvent(ctx context.Context, accessToken string) error {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	accessTokenDetails, err := parseAccessTokenNoValidate(accessToken)
	if err != nil {
		return fmt.Errorf("parse access token: %w", err)
	}

	sessionID, err := idformat.Session.Parse(accessTokenDetails.Session.Id)
	if err != nil {
		return apierror.NewInvalidArgumentError("invalid session ID", fmt.Errorf("parse session ID: %w", err))
	}
	qSession, err := q.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.NewUnauthenticatedError("invalid session", fmt.Errorf("invalid session"))
		}
		return fmt.Errorf("get session: %w", err)
	}

	orgID, err := idformat.Organization.Parse(accessTokenDetails.Organization.Id)
	if err != nil {
		return apierror.NewInvalidArgumentError("invalid organization ID", fmt.Errorf("parse organization ID: %w", err))
	}

	userID, err := idformat.User.Parse(accessTokenDetails.User.Id)
	if err != nil {
		return apierror.NewInvalidArgumentError("invalid user ID", fmt.Errorf("parse user ID: %w", err))
	}

	eventTime := time.Now()
	eventID := uuidv7.NewWithTime(eventTime)

	var impersonatorEmail *string
	if accessTokenDetails.Impersonator != nil && accessTokenDetails.Impersonator.Email != "" {
		impersonatorEmail = &accessTokenDetails.Impersonator.Email
	}

	eventDetails := parseSessionEventDetails(qSession, impersonatorEmail)

	eventDetailsBytes, err := protojson.Marshal(eventDetails)
	if err != nil {
		return fmt.Errorf("marshal event details: %w", err)
	}

	// Since this is being called in a context that doesn't have authn context data,
	// we need to manually set properties like ProjectID and OrganizationID as gleaned
	// from the session details. As such, we can't use the logAuditEvent() fuction
	// directly, so we're calling the CreateAuditLogEvent query directly.
	resourceType := queries.AuditLogEventResourceTypeSession
	if _, err := q.CreateAuditLogEvent(ctx, queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      projectID,
		OrganizationID: (*uuid.UUID)(&orgID),
		UserID:         (*uuid.UUID)(&userID),
		SessionID:      (*uuid.UUID)(&sessionID),
		ResourceType:   &resourceType,
		ResourceID:     (*uuid.UUID)(&sessionID),
		EventName:      "tesseral.sessions.refresh",
		EventTime:      &eventTime,
		EventDetails:   eventDetailsBytes,
	}); err != nil {
		return fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func parseAccessTokenNoValidate(accessToken string) (*commonv1.AccessTokenData, error) {
	jwtParts := strings.Split(accessToken, ".")
	if len(jwtParts) != 3 {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("invalid access token format: expected 3 parts, got %d", len(jwtParts)))
	}

	payloadSegment := jwtParts[1]
	decoded, err := base64.RawURLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to decode access token payload: %w", err))
	}

	var data commonv1.AccessTokenData

	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid credential", fmt.Errorf("failed to unmarshal access token claims: %w", err))
	}

	return &data, nil
}
