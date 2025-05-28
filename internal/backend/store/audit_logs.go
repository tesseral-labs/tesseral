package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/common/auditlog"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) CreateAuditLogEvent(ctx context.Context, req *backendv1.CreateAuditLogEventRequest) (*backendv1.CreateAuditLogEventResponse, error) {
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
	if orgID == "" && userID == "" && sessionID == "" && apiKeyID == "" {
		return nil, apierror.NewInvalidArgumentError("", errors.New("either organization_id, user_id, session_id, or api_key_id must be provided"))
	}
	if sessionID != "" {
		id, err := idformat.Session.Parse(sessionID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", err)
		}
		sessionUUID = (*uuid.UUID)(&id)
	}
	if userID != "" {
		id, err := idformat.User.Parse(userID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", err)
		}
		userUUID = (*uuid.UUID)(&id)
	}
	if apiKeyID != "" {
		id, err := idformat.APIKey.Parse(apiKeyID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", err)
		}
		apiKeyUUID = (*uuid.UUID)(&id)
	}
	if orgID != "" {
		id, err := idformat.Organization.Parse(orgID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", err)
		}
		orgUUID = (*uuid.UUID)(&id)
	}

	// Marshal the details to JSON if provided.
	var eventDetailsJSON []byte
	if eventDetails := req.AuditLogEvent.EventDetails; eventDetails != nil {
		json, err := eventDetails.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("create audit log event: failed to marshal event details JSON: %w", err)
		}
		eventDetailsJSON = json
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

	return &backendv1.CreateAuditLogEventResponse{
		AuditLogEvent: event,
	}, nil
}
