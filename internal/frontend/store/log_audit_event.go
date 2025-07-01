package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
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

func (s *Store) logAuditEvent(ctx context.Context, q *queries.Queries, req logAuditEventParams) (queries.AuditLogEvent, error) {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now()
	eventID := uuidv7.NewWithTime(eventTime)

	eventDetailsBytes, err := protojson.Marshal(req.EventDetails)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to marshal event details: %w", err)
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: refOrNil(authn.OrganizationID(ctx)),
		ActorUserID:    refOrNil(authn.UserID(ctx)),
		ActorSessionID: refOrNil(authn.SessionID(ctx)),
		ResourceType:   refOrNil(req.ResourceType),
		ResourceID:     req.ResourceID,
		EventName:      req.EventName,
		EventTime:      &eventTime,
		EventDetails:   eventDetailsBytes,
	}

	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return queries.AuditLogEvent{}, err
	}

	return qEvent, nil
}

func parseSessionEventDetails(qSession queries.Session, impersonatorEmail *string) *frontendv1.Session {
	var primaryAuthFactor frontendv1.PrimaryAuthFactor
	switch qSession.PrimaryAuthFactor {
	case queries.PrimaryAuthFactorEmail:
		primaryAuthFactor = frontendv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_EMAIL
	case queries.PrimaryAuthFactorGoogle:
		primaryAuthFactor = frontendv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GOOGLE
	case queries.PrimaryAuthFactorMicrosoft:
		primaryAuthFactor = frontendv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_MICROSOFT
	case queries.PrimaryAuthFactorGithub:
		primaryAuthFactor = frontendv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GITHUB
	case queries.PrimaryAuthFactorSaml:
		primaryAuthFactor = frontendv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_SAML
	case queries.PrimaryAuthFactorOidc:
		primaryAuthFactor = frontendv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_OIDC
	default:
		primaryAuthFactor = frontendv1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_UNSPECIFIED
	}

	return &frontendv1.Session{
		Id:                idformat.Session.Format(qSession.ID),
		UserId:            idformat.User.Format(qSession.UserID),
		ExpireTime:        timestamppb.New(derefOrEmpty(qSession.ExpireTime)),
		LastActiveTime:    timestamppb.New(derefOrEmpty(qSession.LastActiveTime)),
		PrimaryAuthFactor: primaryAuthFactor,
		ImpersonatorEmail: derefOrEmpty(impersonatorEmail),
	}
}
