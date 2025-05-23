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
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ActorType string

const (
	ActorTypeUser   ActorType = "user"
	ActorTypeApiKey ActorType = "api_key"
)

func (s *Store) CreateAuditLogEvent(ctx context.Context, req *backendv1.CreateAuditLogEventRequest) (*backendv1.CreateAuditLogEventResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	// TODO: Feature flag check?

	eventTime := time.Now().UTC()
	if req.AuditLogEvent.EventTime != nil {
		eventTime = req.AuditLogEvent.EventTime.AsTime()
	}

	// Generate the UUIDv7 based on the event time.
	id, err := makeUUIDv7(eventTime)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: failed to create UUID: %w", err)
	}
	eventName := req.AuditLogEvent.EventName
	if eventName == "" {
		return nil, apierror.NewInvalidArgumentError("", errors.New("missing event name"))
	}

	// Resolve the actor type/ID from the given inputs.
	var (
		orgID       = req.AuditLogEvent.OrganizationId
		orgUUID     uuid.UUID
		userID      = req.AuditLogEvent.GetUserId()
		userUUID    *uuid.UUID
		sessionID   = req.AuditLogEvent.GetSessionId()
		sessionUUID *uuid.UUID
		apiKeyID    = req.AuditLogEvent.GetApiKeyId()
		apiKeyUUID  *uuid.UUID
	)
	switch {
	case sessionID != "":
		id, err := idformat.Session.Parse(sessionID)
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
		orgUUID = user.OrganizationID
	case userID != "":
		id, err := idformat.User.Parse(userID)
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
		orgUUID = user.OrganizationID
	case apiKeyID != "":
		id, err := idformat.APIKey.Parse(apiKeyID)
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
		orgUUID = apiKey.OrganizationID
	case orgID != "":
		orgUUID, err = idformat.Organization.Parse(orgID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", err)
		}
		_, err = q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ID:        orgUUID,
			ProjectID: projectID,
		})
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization_id", err)
		}
	default:
		return nil, apierror.NewInvalidArgumentError("", errors.New("either organization_id, user_id, session_id, or api_key_id must be provided"))
	}

	// Marshal the details to JSON if provided.
	var detailsJSON []byte
	if details := req.AuditLogEvent.EventDetails; details != nil {
		detailsJSON, err = details.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("create audit log event: failed to marshal event details JSON: %w", err)
		}
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             id,
		OrganizationID: orgUUID,
		UserID:         userUUID,
		SessionID:      sessionUUID,
		ApiKeyID:       apiKeyUUID,
		EventName:      eventName,
		EventDetails:   detailsJSON,
	}
	qEvent, err := q.CreateAuditLogEvent(ctx, qEventParams)
	if err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	event, err := parseAuditLogEvent(qEvent)
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

// makeUUIDv7 copies google/uuid for constructing a UUIDv7 at a point in time
// (as opposed to making one for the current instant).
func makeUUIDv7(ts time.Time) (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, err
	}

	nano := ts.UnixNano()
	const nanoPerMilli = 1_000_000
	milli := nano / nanoPerMilli

	// Sequence number is not used since there is no accurate way to establish one.
	// Instead we leave the random bits from the V4 in place.

	/*
		 0                   1                   2                   3
		 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                           unix_ts_ms                          |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|          unix_ts_ms           |  ver  |  rand_a (12 bit seq)  |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|var|                        rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                            rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
	_ = id[15] // bounds check

	id[0] = byte(milli >> 40)
	id[1] = byte(milli >> 32)
	id[2] = byte(milli >> 24)
	id[3] = byte(milli >> 16)
	id[4] = byte(milli >> 8)
	id[5] = byte(milli)

	// Swap version for V7
	id[6] = 0x70 | (0x0F & id[6])

	return id, nil
}

// tsFromUUIDv7 reconstructs a timestamp from a UUIDv7, accurate to millisecond precision.
func tsFromUUIDv7(id uuid.UUID) time.Time {
	milli := (int64(id[0]) << 40) | (int64(id[1]) << 32) | (int64(id[2]) << 24) | (int64(id[3]) << 16) | (int64(id[4]) << 8) | int64(id[5])
	return time.UnixMilli(milli).UTC()
}

func parseAuditLogEvent(qEvent queries.OrganizationAuditLogEvent) (*backendv1.AuditLogEvent, error) {
	detailsJSON := qEvent.EventDetails
	var details structpb.Struct
	if err := details.UnmarshalJSON(detailsJSON); err != nil {
		return nil, err
	}

	ts := tsFromUUIDv7(qEvent.ID)

	var (
		userID    *string
		sessionID *string
		apiKeyID  *string
	)
	if userUUID := qEvent.UserID; userUUID != nil {
		id := idformat.User.Format(*userUUID)
		userID = &id
	}
	if sessionUUID := qEvent.SessionID; sessionUUID != nil {
		id := idformat.Session.Format(*sessionUUID)
		sessionID = &id
	}
	if apiKeyUUID := qEvent.ApiKeyID; apiKeyUUID != nil {
		id := idformat.APIKey.Format(*apiKeyUUID)
		apiKeyID = &id
	}

	return &backendv1.AuditLogEvent{
		Id:             idformat.AuditLogEvent.Format(qEvent.ID),
		OrganizationId: idformat.Organization.Format(qEvent.OrganizationID),
		UserId:         userID,
		SessionId:      sessionID,
		ApiKeyId:       apiKeyID,
		EventName:      qEvent.EventName,
		EventTime:      timestamppb.New(ts),
		EventDetails:   &details,
	}, nil
}
