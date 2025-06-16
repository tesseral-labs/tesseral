package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type logAuditEventParams struct {
	OrganizationID *uuid.UUID

	EventName    string
	EventDetails proto.Message
	ResourceType queries.AuditLogEventResourceType
	ResourceID   *uuid.UUID
}

func (s *Store) logAuditEvent(ctx context.Context, q *queries.Queries, data logAuditEventParams) (queries.AuditLogEvent, error) {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now()
	eventID := uuidv7.NewWithTime(eventTime)

	eventDetailsBytes, err := protojson.Marshal(data.EventDetails)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to marshal event details: %w", err)
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: data.OrganizationID,
		ResourceType:   refOrNil(data.ResourceType),
		ResourceID:     data.ResourceID,
		EventName:      data.EventName,
		EventTime:      &eventTime,
		EventDetails:   eventDetailsBytes,
	}

	contextData := authn.GetContextData(ctx)
	switch {
	case contextData.ProjectAPIKey != nil:
		backendApiKeyUUID, err := idformat.BackendAPIKey.Parse(contextData.ProjectAPIKey.BackendAPIKeyID)
		if err != nil {
			return queries.AuditLogEvent{}, fmt.Errorf("parse backend api key id: %w", err)
		}
		qEventParams.ActorBackendApiKeyID = (*uuid.UUID)(&backendApiKeyUUID)
	case contextData.DogfoodSession != nil:
		dogfoodUserUUID, err := idformat.User.Parse(contextData.DogfoodSession.UserID)
		if err != nil {
			return queries.AuditLogEvent{}, fmt.Errorf("parse dogfood session user id: %w", err)
		}
		qEventParams.ActorConsoleUserID = (*uuid.UUID)(&dogfoodUserUUID)

		dogfoodSessionUUID, err := idformat.Session.Parse(contextData.DogfoodSession.SessionID)
		if err != nil {
			return queries.AuditLogEvent{}, fmt.Errorf("parse dogfood session project id: %w", err)
		}
		qEventParams.ActorSessionID = (*uuid.UUID)(&dogfoodSessionUUID)
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	return qEvent, nil
}
