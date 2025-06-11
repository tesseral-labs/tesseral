package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/types/known/structpb"
)

type logAuditEventParams struct {
	EventName      string
	EventDetails   map[string]any
	OrganizationID *uuid.UUID
	ResourceType   queries.AuditLogEventResourceType
	ResourceID     *uuid.UUID
}

func (s *Store) logAuditEvent(ctx context.Context, q *queries.Queries, data logAuditEventParams) (queries.AuditLogEvent, error) {
	// Generate the UUIDv7 based on the event time.
	eventTime := time.Now()
	eventID, err := uuidv7.NewWithTime(eventTime)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to create UUID: %w", err)
	}

	eventDetailsBytes, err := json.Marshal(data.EventDetails)
	if err != nil {
		return queries.AuditLogEvent{}, fmt.Errorf("failed to marshal event details: %w", err)
	}

	qEventParams := queries.CreateAuditLogEventParams{
		ID:             eventID,
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: refOrNil(authn.OrganizationID(ctx)),
		UserID:         refOrNil(authn.UserID(ctx)),
		SessionID:      refOrNil(authn.SessionID(ctx)),
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

func (s *Store) ListAuditLogEvents(ctx context.Context, req *frontendv1.ListAuditLogEventsRequest) (*frontendv1.ListAuditLogEventsResponse, error) {
	tx, _, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)
	orgID := authn.OrganizationID(ctx)

	if err := s.validateIsOwner(ctx); err != nil {
		return nil, err
	}

	filter := new(frontendv1.ListAuditLogEventsRequest)
	if req.PageToken != "" {
		if err := s.pageEncoder.Unmarshal(req.PageToken, filter); err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid page_token", err)
		}
	} else {
		filter = req
	}

	var startTime time.Time
	if filterStartTime := filter.GetFilterStartTime(); filterStartTime != nil {
		startTime = filterStartTime.AsTime()
	}
	startID := uuidv7.NilWithTime(startTime)

	var endTime time.Time
	if filterEndTime := filter.GetFilterEndTime(); filterEndTime != nil {
		endTime = filterEndTime.AsTime()
	} else {
		endTime = time.Now()
	}
	if endTime.Before(startTime) {
		return nil, apierror.NewInvalidArgumentError("end_time must be after start_time", fmt.Errorf("end time %s is before start time %s", endTime, startTime))
	}
	endID := uuidv7.NilWithTime(endTime)

	wheres := []sq.Sqlizer{
		sq.Eq{"project_id": projectID[:]},
		sq.Eq{"organization_id": orgID[:]},
		sq.Lt{"id": endID},
	}
	if !startTime.IsZero() {
		wheres = append(wheres, sq.Gt{"id": startID})
	}
	if len(filter.GetFilterEventName()) > 0 {
		wheres = append(wheres, sq.Eq{"event_name": filter.FilterEventName})
	}
	if userID := filter.GetFilterUserId(); userID != "" {
		userID, err := idformat.User.Parse(userID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.user_id", err)
		}
		wheres = append(wheres, sq.Eq{"user_id": userID[:]})
	}
	if sessionID := filter.GetFilterSessionId(); sessionID != "" {
		sessionID, err := idformat.Session.Parse(sessionID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.session_id", err)
		}
		wheres = append(wheres, sq.Eq{"session_id": sessionID[:]})
	}
	if apiKeyID := filter.GetFilterApiKeyId(); apiKeyID != "" {
		apiKeyID, err := idformat.APIKey.Parse(apiKeyID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.api_key_id", err)
		}
		wheres = append(wheres, sq.Eq{"api_key_id": apiKeyID[:]})
	}

	const limit = 10

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("*").
		From("audit_log_events").
		Where(sq.And(wheres)).
		OrderBy("id desc").
		Limit(limit + 1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to construct sql query: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}
	results, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*frontendv1.AuditLogEvent, error) {
		var dto queries.AuditLogEvent
		if err := row.Scan(
			&dto.ID,
			&dto.ProjectID,
			&dto.OrganizationID,
			&dto.UserID,
			&dto.SessionID,
			&dto.ApiKeyID,
			&dto.DogfoodUserID,
			&dto.DogfoodSessionID,
			&dto.BackendApiKeyID,
			&dto.IntermediateSessionID,
			&dto.ResourceType,
			&dto.ResourceID,
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
	if len(results) == limit+1 {
		last := results[limit-1]
		filter.FilterEndTime = last.EventTime
		nextPageToken = s.pageEncoder.Marshal(filter)
		results = results[:limit]
	}

	return &frontendv1.ListAuditLogEventsResponse{
		AuditLogEvents: results,
		NextPageToken:  nextPageToken,
	}, nil
}

func parseAuditLogEvent(qEvent queries.AuditLogEvent) (*frontendv1.AuditLogEvent, error) {
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

	return &frontendv1.AuditLogEvent{
		Id:           idformat.AuditLogEvent.Format(qEvent.ID),
		UserId:       userID,
		SessionId:    sessionID,
		ApiKeyId:     apiKeyID,
		EventName:    qEvent.EventName,
		EventTime:    timestampOrNil(qEvent.EventTime),
		EventDetails: &eventDetails,
	}, nil
}

type sessionEventDetails struct {
	ID                string                       `json:"id"`
	UserID            string                       `json:"userId"`
	ExpireTime        string                       `json:"expireTime,omitempty"`
	LastActiveTime    string                       `json:"lastActiveTime,omitempty"`
	PrimaryAuthFactor frontendv1.PrimaryAuthFactor `json:"primaryAuthFactor,omitempty"`
	ImpersonatorEmail string                       `json:"impersonatorEmail,omitempty"`
}

func parseSessionEventDetails(qSession queries.Session, impersonatorEmail *string) *sessionEventDetails {
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

	return &sessionEventDetails{
		ID:                idformat.Session.Format(qSession.ID),
		UserID:            idformat.User.Format(qSession.UserID),
		ExpireTime:        qSession.ExpireTime.String(),
		LastActiveTime:    qSession.LastActiveTime.String(),
		PrimaryAuthFactor: primaryAuthFactor,
		ImpersonatorEmail: derefOrEmpty(impersonatorEmail),
	}
}
