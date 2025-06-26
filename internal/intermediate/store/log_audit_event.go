package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type logAuditEventParams struct {
	EventName      string
	EventDetails   proto.Message
	OrganizationID *uuid.UUID
	ResourceType   queries.AuditLogEventResourceType
	ResourceID     *uuid.UUID
}

func (s *Store) logAuditEvent(ctx context.Context, q *queries.Queries, data logAuditEventParams) (queries.AuditLogEvent, error) {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now()
	eventID := uuidv7.NewWithTime(eventTime)

	eventDetailsBytes, err := protojson.Marshal(data.EventDetails)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to marshal event details: %w", err)
	}

	intermediateSessionID := authn.IntermediateSessionID(ctx)
	qEventParams := queries.CreateAuditLogEventParams{
		ID:                         eventID,
		ProjectID:                  authn.ProjectID(ctx),
		OrganizationID:             data.OrganizationID,
		ResourceType:               &data.ResourceType,
		ResourceID:                 data.ResourceID,
		EventName:                  data.EventName,
		EventTime:                  &eventTime,
		EventDetails:               eventDetailsBytes,
		ActorIntermediateSessionID: &intermediateSessionID,
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	return qEvent, nil
}
