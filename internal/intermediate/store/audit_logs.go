package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func parseSessionEventDetails(qSession queries.Session, impersonatorEmail *string) *intermediatev1.SessionCreated {
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

	return &intermediatev1.SessionCreated{
		Session: &intermediatev1.Session{
			Id:                idformat.Session.Format(qSession.ID),
			UserId:            idformat.User.Format(qSession.UserID),
			ExpireTime:        timestamppb.New(derefOrEmpty(qSession.ExpireTime)),
			LastActiveTime:    timestamppb.New(derefOrEmpty(qSession.LastActiveTime)),
			PrimaryAuthFactor: primaryAuthFactor,
			ImpersonatorEmail: derefOrEmpty(impersonatorEmail),
		},
	}
}
