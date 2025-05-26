package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/structpb"
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
	var eventDetailsJSON []byte
	if eventDetails := req.AuditLogEvent.EventDetails; eventDetails != nil {
		eventDetailsJSON, err = eventDetails.MarshalJSON()
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
		EventTime:      &eventTime,
		EventDetails:   eventDetailsJSON,
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
	return makeUUIDv7Base(ts, id), nil
}

// makeUUIDv7 copies google/uuid for constructing a UUIDv7 at a point in time
// given a base UUID data.
func makeUUIDv7Base(ts time.Time, id uuid.UUID) uuid.UUID {
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

	id[6] = 0x70 | (0x0F & id[6]) // Version is 7 (0b0111)
	id[8] = 0x80 | (0x3F & id[8]) // Variant is 0b10

	return id
}

// tsFromUUIDv7 reconstructs a timestamp from a UUIDv7, accurate to millisecond precision.
func tsFromUUIDv7(id uuid.UUID) time.Time {
	milli := (int64(id[0]) << 40) | (int64(id[1]) << 32) | (int64(id[2]) << 24) | (int64(id[3]) << 16) | (int64(id[4]) << 8) | int64(id[5])
	return time.UnixMilli(milli).UTC()
}

func parseAuditLogEvent(qEvent queries.OrganizationAuditLogEvent) (*backendv1.AuditLogEvent, error) {
	eventDetailsJSON := qEvent.EventDetails
	var eventDetails structpb.Struct
	if err := eventDetails.UnmarshalJSON(eventDetailsJSON); err != nil {
		return nil, err
	}

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
		EventTime:      timestampOrNil(qEvent.EventTime),
		EventDetails:   &eventDetails,
	}, nil
}

func (s *Store) ListAuditLogEvents(ctx context.Context, req *backendv1.ListAuditLogEventsRequest) (*backendv1.ListAuditLogEventsResponse, error) {
	tx, _, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	var filter *backendv1.ListAuditLogEventsRequest_Filter
	if req.PageToken != "" {
		filter = new(backendv1.ListAuditLogEventsRequest_Filter)
		if err := s.pageEncoder.Unmarshal(req.PageToken, filter); err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid page_token", err)
		}
	} else {
		filter = req.Filter
	}

	limit := uint64(10)
	if req.PageSize != 0 {
		limit = uint64(req.PageSize)
	}

	startTime := filter.GetStartTime().AsTime()
	startID := makeUUIDv7Base(startTime, uuid.UUID{})

	wheres := []sq.Sqlizer{
		sq.Eq{"organization_id": orgID[:]},
		sq.Gt{"id": startID},
	}
	if len(filter.GetEventName()) > 0 {
		wheres = append(wheres, sq.Eq{"event_name": filter.EventName})
	}
	if filter.GetUserId() != "" {
		userID, err := idformat.User.Parse(filter.UserId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.user_id", err)
		}
		wheres = append(wheres, sq.Eq{"user_id": userID[:]})
	}
	if filter.GetSessionId() != "" {
		sessionID, err := idformat.Session.Parse(filter.SessionId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.session_id", err)
		}
		wheres = append(wheres, sq.Eq{"session_id": sessionID[:]})
	}
	if filter.GetApiKeyId() != "" {
		apiKeyID, err := idformat.APIKey.Parse(filter.ApiKeyId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.api_key_id", err)
		}
		wheres = append(wheres, sq.Eq{"api_key_id": apiKeyID[:]})
	}
	orderBy := []string{}
	if req.OrderBy != "" {
		orderBy = append(orderBy, req.OrderBy)
	}
	orderBy = append(orderBy, "id desc")

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("organization_audit_log_events").
		Where(sq.And(wheres)).
		OrderBy(orderBy...).
		Limit(limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to construct sql query: %w", err)
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}
	results, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*backendv1.AuditLogEvent, error) {
		var dto queries.OrganizationAuditLogEvent
		if err := row.Scan(
			&dto.ID,
			&dto.OrganizationID,
			&dto.UserID,
			&dto.SessionID,
			&dto.ApiKeyID,
			&dto.EventName,
			&dto.EventTime,
			&dto.EventDetails,
		); err != nil {
			return nil, err
		}
		return parseAuditLogEvent(dto)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to collect audit log events: %w", err)
	}

	var nextPageToken string
	if len(results) > 0 {
		last := results[len(results)-1]
		filter.StartTime = last.EventTime
		nextPageToken = s.pageEncoder.Marshal(filter)
	}

	return &backendv1.ListAuditLogEventsResponse{
		AuditLogEvents: results,
		NextPageToken:  nextPageToken,
	}, nil
}
