package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/common/auditlog"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/common/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	var (
		projectID = event.ProjectID
		orgID     = event.OrganizationID
		userID    = event.UserID
		sessionID = event.SessionID
		apiKeyID  = event.ApiKeyID
	)
	switch {
	case projectID == uuid.Nil:
		return nil, errors.New("missing project_id")
	case sessionID != nil:
		// Lookup session, user, and organization
		session, err := q.GetSession(ctx, queries.GetSessionParams{
			ID:        *sessionID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("get session: %w", err))
		}

		user, err := q.GetUser(ctx, queries.GetUserParams{
			ID:        session.UserID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid session_id", fmt.Errorf("get user: %w", err))
		}
		userID = (*uuid.UUID)(&user.ID)
		orgID = (*uuid.UUID)(&user.OrganizationID)
	case userID != nil:
		// Lookup user and organization
		user, err := q.GetUser(ctx, queries.GetUserParams{
			ID:        *userID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid user_id", fmt.Errorf("get user: %w", err))
		}
		orgID = (*uuid.UUID)(&user.OrganizationID)
	case apiKeyID != nil:
		// Lookup API key and organization
		apiKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
			ID:        *apiKeyID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid api_key_id", fmt.Errorf("get api key: %w", err))
		}
		orgID = (*uuid.UUID)(&apiKey.OrganizationID)
	case orgID != nil:
		_, err = q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ID:        *orgID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", fmt.Errorf("get organization: %w", err))
		}
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             event.ID,
		ProjectID:      projectID,
		OrganizationID: orgID,
		UserID:         userID,
		SessionID:      sessionID,
		ApiKeyID:       apiKeyID,
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
		organizationID string
		userID         *string
		sessionID      *string
		apiKeyID       *string
	)
	if orgUUID := qEvent.OrganizationID; orgUUID != nil {
		organizationID = idformat.Organization.Format(*orgUUID)
	}
	if userUUID := qEvent.UserID; userUUID != nil {
		userID_ := idformat.User.Format(*userUUID)
		userID = &userID_
	}
	if sessionUUID := qEvent.SessionID; sessionUUID != nil {
		sessionID_ := idformat.Session.Format(*sessionUUID)
		sessionID = &sessionID_
	}
	if apiKeyUUID := qEvent.ApiKeyID; apiKeyUUID != nil {
		apiKeyID_ := idformat.APIKey.Format(*apiKeyUUID)
		apiKeyID = &apiKeyID_
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
