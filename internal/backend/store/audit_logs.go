package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/common/auditlog"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Store) CreateAuditLogEvent(ctx context.Context, req *backendv1.CreateAuditLogEventRequest) (*backendv1.CreateAuditLogEventResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	eventTime := time.Now().UTC()
	if req.AuditLogEvent.EventTime != nil {
		eventTime = req.AuditLogEvent.EventTime.AsTime()
	}
	eventName := req.AuditLogEvent.EventName
	if eventName == "" {
		return nil, apierror.NewInvalidArgumentError("", errors.New("missing event name"))
	}

	// Resolve the actor type/ID from the given inputs.
	var (
		orgID       = req.AuditLogEvent.GetOrganizationId()
		orgUUID     *uuid.UUID
		userID      = req.AuditLogEvent.GetUserId()
		userUUID    *uuid.UUID
		sessionID   = req.AuditLogEvent.GetSessionId()
		sessionUUID *uuid.UUID
		apiKeyID    = req.AuditLogEvent.GetApiKeyId()
		apiKeyUUID  *uuid.UUID
	)
	switch {
	case sessionID != nil:
		id, err := idformat.Session.Parse(sessionID.Value)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", err)
		}
		// Lookup session, user, and organization
		session, err := q.GetSession(ctx, queries.GetSessionParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", err)
		}
		sessionUUID = (*uuid.UUID)(&session.ID)

		user, err := q.GetUser(ctx, queries.GetUserParams{
			ID:        session.UserID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("get user: %w", err))
		}
		userUUID = (*uuid.UUID)(&user.ID)
		orgUUID = (*uuid.UUID)(&user.OrganizationID)
	case userID != nil:
		id, err := idformat.User.Parse(userID.Value)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", err)
		}
		// Lookup user and organization
		user, err := q.GetUser(ctx, queries.GetUserParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", err)
		}
		userUUID = (*uuid.UUID)(&user.ID)
		orgUUID = (*uuid.UUID)(&user.OrganizationID)
	case apiKeyID != nil:
		id, err := idformat.APIKey.Parse(apiKeyID.Value)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", err)
		}
		// Lookup API key and organization
		apiKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", err)
		}
		apiKeyUUID = (*uuid.UUID)(&apiKey.ID)
		orgUUID = (*uuid.UUID)(&apiKey.OrganizationID)
	case orgID != nil:
		id, err := idformat.Organization.Parse(orgID.Value)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", err)
		}
		_, err = q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ID:        id,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", err)
		}
		orgUUID = (*uuid.UUID)(&id)
	default:
		return nil, apierror.NewInvalidArgumentError("", errors.New("either organization_id, user_id, session_id, or api_key_id must be provided"))
	}

	// Marshal the details to JSON if provided.
	var eventDetailsJSON []byte
	if eventDetails := req.AuditLogEvent.EventDetails; eventDetails != nil {
		eventDetailsJSON, err = eventDetails.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("create audit log event: failed to marshal event details JSON: %w", err)
		}
	}

	event, err := s.common.CreateAuditLogEvent(ctx, auditlog.Event{
		ProjectID:      projectID,
		OrganizationID: orgUUID,
		UserID:         userUUID,
		SessionID:      sessionUUID,
		ApiKeyID:       apiKeyUUID,
		EventName:      eventName,
		EventTime:      &eventTime,
		EventDetails:   eventDetailsJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	return &backendv1.CreateAuditLogEventResponse{
		AuditLogEvent: event,
	}, nil
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
