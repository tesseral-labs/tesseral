package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/muststructpb"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type logAuditEventParams struct {
	EventName      string
	EventDetails   *structpb.Value
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
		ID:                    eventID,
		ProjectID:             authn.ProjectID(ctx),
		OrganizationID:        data.OrganizationID,
		ResourceType:          &data.ResourceType,
		ResourceID:            data.ResourceID,
		EventName:             data.EventName,
		EventTime:             &eventTime,
		EventDetails:          eventDetailsBytes,
		IntermediateSessionID: &intermediateSessionID,
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	return qEvent, nil
}

func parseSessionEventDetails(qSession queries.Session, impersonatorEmail *string) *structpb.Value {
	var primaryAuthFactor string
	switch qSession.PrimaryAuthFactor {
	case queries.PrimaryAuthFactorEmail:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_EMAIL.String()
	case queries.PrimaryAuthFactorGoogle:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GOOGLE.String()
	case queries.PrimaryAuthFactorMicrosoft:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_MICROSOFT.String()
	case queries.PrimaryAuthFactorGithub:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GITHUB.String()
	default:
		primaryAuthFactor = intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_UNSPECIFIED.String()
	}

	sessionDetails := map[string]interface{}{
		"session": map[string]interface{}{
			"id":                  idformat.Session.Format(qSession.ID),
			"user_id":             idformat.User.Format(qSession.UserID),
			"expire_time":         qSession.ExpireTime.String(),
			"last_active_time":    qSession.LastActiveTime.String(),
			"primary_auth_factor": primaryAuthFactor,
			"impersonator_email":  derefOrEmpty(impersonatorEmail),
		},
	}

	return structpb.NewStructValue(muststructpb.MustNewStruct(sessionDetails))
}
