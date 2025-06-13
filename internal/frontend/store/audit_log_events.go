package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListAuditLogEvents(ctx context.Context, req *frontendv1.ListAuditLogEventsRequest) (*frontendv1.ListAuditLogEventsResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, err
	}

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// We want data sorted by event time, newest first. That corresponds to
	// paginating through IDs high-to-low, because IDs are uuidv7s for this
	// table.
	startID := uuid.Max
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	listParams := queries.ListAuditLogEventsParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: refOrNil(authn.OrganizationID(ctx)),
		ID:             startID,
		Limit:          int32(limit + 1),
	}

	if req.FilterStartTime != nil {
		filterStartTime := req.FilterStartTime.AsTime()
		listParams.StartTime = &filterStartTime
	}

	if req.FilterEndTime != nil {
		endTime := req.FilterEndTime.AsTime()
		listParams.EndTime = &endTime
	}

	if req.FilterEventName != "" {
		listParams.EventName = &req.FilterEventName
	}

	if req.FilterUserId != "" {
		id, err := idformat.User.Parse(req.FilterUserId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", err)
		}
		listParams.ActorUserID = (*uuid.UUID)(&id)
	}

	qAuditLogEvents, err := q.ListAuditLogEvents(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("list audit log events: %w", err)
	}

	var auditLogEvents []*frontendv1.AuditLogEvent
	for _, qAuditLogEvent := range qAuditLogEvents {
		event, err := parseAuditLogEvent(qAuditLogEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse audit log event: %w", err)
		}
		auditLogEvents = append(auditLogEvents, event)
	}

	var nextPageToken string
	if len(qAuditLogEvents) == limit+1 {
		last := qAuditLogEvents[limit-1]
		nextPageToken = s.pageEncoder.Marshal(last.EventTime)
		auditLogEvents = auditLogEvents[:limit]
	}

	return &frontendv1.ListAuditLogEventsResponse{
		AuditLogEvents: auditLogEvents,
		NextPageToken:  nextPageToken,
	}, nil
}

func parseAuditLogEvent(qAuditLogEvent queries.AuditLogEvent) (*frontendv1.AuditLogEvent, error) {
	var eventDetails structpb.Struct
	if err := protojson.Unmarshal(qAuditLogEvent.EventDetails, &eventDetails); err != nil {
		return nil, fmt.Errorf("unmarshal event details: %w", err)
	}

	var actorUserID string
	if qAuditLogEvent.ActorUserID != nil {
		actorUserID = idformat.User.Format(*qAuditLogEvent.ActorUserID)
	}

	var actorSessionID string
	if qAuditLogEvent.ActorSessionID != nil {
		actorSessionID = idformat.Session.Format(*qAuditLogEvent.ActorSessionID)
	}

	var actorAPIKeyID string
	if qAuditLogEvent.ActorApiKeyID != nil {
		actorAPIKeyID = idformat.APIKey.Format(*qAuditLogEvent.ActorApiKeyID)
	}

	var actorIntermediateSessionID string
	if qAuditLogEvent.ActorIntermediateSessionID != nil {
		actorIntermediateSessionID = idformat.IntermediateSession.Format(*qAuditLogEvent.ActorIntermediateSessionID)
	}

	return &frontendv1.AuditLogEvent{
		Id:                         idformat.AuditLogEvent.Format(qAuditLogEvent.ID),
		ActorUserId:                actorUserID,
		ActorSessionId:             actorSessionID,
		ActorApiKeyId:              actorAPIKeyID,
		ActorIntermediateSessionId: actorIntermediateSessionID,
		EventName:                  qAuditLogEvent.EventName,
		EventTime:                  timestampOrNil(qAuditLogEvent.EventTime),
		EventDetails:               &eventDetails,
	}, nil
}

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

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: refOrNil(authn.OrganizationID(ctx)),
		ActorUserID:    refOrNil(authn.UserID(ctx)),
		ActorSessionID: refOrNil(authn.SessionID(ctx)),
		ResourceType:   refOrNil(data.ResourceType),
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
