package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
)

type logAuditEventParams struct {
	EventName      string
	EventDetails   map[string]any
	OrganizationID *uuid.UUID
	ResourceType   queries.AuditLogEventResourceType
	ResourceID     *uuid.UUID
}

func (s *Store) logAuditEvent(ctx context.Context, q *queries.Queries, data logAuditEventParams) (queries.AuditLogEvent, error) {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now().UTC()
	eventID, err := uuidv7.NewWithTime(eventTime)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to create UUID: %w", err)
	}

	eventDetailsBytes, err := json.Marshal(data.EventDetails)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to marshal event details: %w", err)
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

	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	return qEvent, nil
}

// The intermediate service doesn't do CRUD in the traditional sense,
// so we use a custom struct to represent the event details for session events.
// This struct is used to format the session event details for logging.
type sessionEventDetails struct {
	ID                string                           `json:"id"`
	UserID            string                           `json:"userId"`
	ExpireTime        string                           `json:"expireTime,omitempty"`
	LastActiveTime    string                           `json:"lastActiveTime,omitempty"`
	PrimaryAuthFactor intermediatev1.PrimaryAuthFactor `json:"primaryAuthFactor,omitempty"`
	ImpersonatorEmail string                           `json:"impersonatorEmail,omitempty"`
}

func parseSessionEventDetails(qSession queries.Session, impersonatorEmail *string) *sessionEventDetails {
	var primaryAuthFactor intermediatev1.PrimaryAuthFactor
	switch qSession.PrimaryAuthFactor {
	case queries.PrimaryAuthFactorEmail:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_EMAIL
	case queries.PrimaryAuthFactorGoogle:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GOOGLE
	case queries.PrimaryAuthFactorMicrosoft:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_MICROSOFT
	case queries.PrimaryAuthFactorGithub:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GITHUB
	default:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_UNSPECIFIED
	}

	return &sessionEventDetails{
		ID:                idformat.Session.Format(qSession.ID),
		UserID:            idformat.User.Format(qSession.UserID),
		ExpireTime:        qSession.ExpireTime.String(),
		LastActiveTime:    qSession.LastActiveTime.String(),
		PrimaryAuthFactor: primaryAuthFactor,
		ImpersonatorEmail: derefOrEmpty(impersonatorEmail),
	}
}
