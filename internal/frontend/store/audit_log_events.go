package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
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
		nextPageToken = s.pageEncoder.Marshal(qAuditLogEvents[limit].ID)
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
