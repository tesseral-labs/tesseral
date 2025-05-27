package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/auditlog"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Store) CreateAuditLogEvent(ctx context.Context, event auditlog.Event) (*commonv1.AuditLogEvent, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}
	defer rollback()

	eventTime := event.EventTime
	if eventTime == nil {
		now := time.Now().UTC()
		eventTime = &now
	}

	// Generate the UUIDv7 based on the event time.
	if event.ID == uuid.Nil {
		event.ID, err = uuidv7.NewWithTime(*eventTime)
		if err != nil {
			return nil, fmt.Errorf("create audit log event: failed to create UUID: %w", err)
		}
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             event.ID,
		ProjectID:      event.ProjectID,
		OrganizationID: event.OrganizationID,
		UserID:         event.UserID,
		SessionID:      event.SessionID,
		ApiKeyID:       event.ApiKeyID,
		EventName:      event.EventName,
		EventTime:      eventTime,
		EventDetails:   event.EventDetails,
	}
	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	pEvent, err := parseAuditLogEvent(qEvent)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	return pEvent, nil
}

func parseAuditLogEvent(qEvent queries.AuditLogEvent) (*commonv1.AuditLogEvent, error) {
	eventDetailsJSON := qEvent.EventDetails
	var eventDetails structpb.Struct
	if err := eventDetails.UnmarshalJSON(eventDetailsJSON); err != nil {
		return nil, err
	}

	var (
		organizationID *wrapperspb.StringValue
		userID         *wrapperspb.StringValue
		sessionID      *wrapperspb.StringValue
		apiKeyID       *wrapperspb.StringValue
	)
	if orgUUID := qEvent.OrganizationID; orgUUID != nil {
		organizationID = wrapperspb.String(idformat.Organization.Format(*orgUUID))
	}
	if userUUID := qEvent.UserID; userUUID != nil {
		userID = wrapperspb.String(idformat.User.Format(*userUUID))
	}
	if sessionUUID := qEvent.SessionID; sessionUUID != nil {
		sessionID = wrapperspb.String(idformat.Session.Format(*sessionUUID))
	}
	if apiKeyUUID := qEvent.ApiKeyID; apiKeyUUID != nil {
		apiKeyID = wrapperspb.String(idformat.APIKey.Format(*apiKeyUUID))
	}

	return &commonv1.AuditLogEvent{
		Id:             idformat.AuditLogEvent.Format(qEvent.ID),
		OrganizationId: organizationID,
		UserId:         userID,
		SessionId:      sessionID,
		ApiKeyId:       apiKeyID,
		EventName:      qEvent.EventName,
		EventTime:      timestampOrNil(qEvent.EventTime),
		EventDetails:   &eventDetails,
	}, nil
}

func timestampOrNil(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}
