package store

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/uuidv7"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Store) ListAuditLogEvents(ctx context.Context, req *frontendv1.ListAuditLogEventsRequest) (*frontendv1.ListAuditLogEventsResponse, error) {
	tx, _, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)
	orgID := authn.OrganizationID(ctx)

	// TODO: Enforce owner-only

	filter := new(frontendv1.ListAuditLogEventsRequest_Filter)
	if req.PageToken != "" {
		if err := s.pageEncoder.Unmarshal(req.PageToken, filter); err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid page_token", err)
		}
	} else if req.Filter != nil {
		filter = req.Filter
	}

	limit := uint64(10)
	if req.PageSize != 0 {
		limit = uint64(req.PageSize)
	}

	var startTime time.Time
	if filterStartTime := filter.GetStartTime(); filterStartTime != nil {
		startTime = filterStartTime.AsTime()
	}
	startID := uuidv7.NilWithTime(startTime)

	var endTime time.Time
	if filterEndTime := filter.GetEndTime(); filterEndTime != nil {
		endTime = filterEndTime.AsTime()
	} else {
		endTime = time.Now().UTC()
	}
	if endTime.Before(startTime) {
		return nil, apierror.NewInvalidArgumentError("end_time must be after start_time", fmt.Errorf("end time %s is before start time %s", endTime, startTime))
	}
	endID := uuidv7.NilWithTime(endTime)

	slog.InfoContext(ctx, "list_audit_log_events", "startTime", startTime, "endTime", endTime, "eventName", filter.EventName, "userId", filter.UserId, "sessionId", filter.SessionId, "apiKeyId", filter.ApiKeyId)

	wheres := []sq.Sqlizer{
		sq.Eq{"project_id": projectID[:]},
		sq.Eq{"organization_id": orgID[:]},
		sq.Lt{"id": endID},
	}
	if !startTime.IsZero() {
		wheres = append(wheres, sq.Gt{"id": startID})
	}
	if len(filter.GetEventName()) > 0 {
		wheres = append(wheres, sq.Eq{"event_name": filter.EventName})
	}
	if userID := filter.GetUserId(); userID != "" {
		userID, err := idformat.User.Parse(userID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.user_id", err)
		}
		wheres = append(wheres, sq.Eq{"user_id": userID[:]})
	}
	if sessionID := filter.GetSessionId(); sessionID != "" {
		sessionID, err := idformat.Session.Parse(sessionID)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid filter.session_id", err)
		}
		wheres = append(wheres, sq.Eq{"session_id": sessionID[:]})
	}
	if apiKeyID := filter.GetApiKeyId(); apiKeyID != "" {
		apiKeyID, err := idformat.APIKey.Parse(apiKeyID)
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
		From("audit_log_events").
		Where(sq.And(wheres)).
		OrderBy(orderBy...).
		Limit(limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to construct sql query: %w", err)
	}

	slog.InfoContext(ctx, "execute_query", "sql", sql, "args", args)
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}
	results, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*commonv1.AuditLogEvent, error) {
		var dto queries.AuditLogEvent
		if err := row.Scan(
			&dto.ID,
			&dto.ProjectID,
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
	if len(results) == int(limit) {
		last := results[len(results)-1]
		if slices.Contains(orderBy, "id asc") {
			filter.StartTime = last.EventTime
		} else {
			filter.EndTime = last.EventTime
		}
		slog.InfoContext(ctx, "next_page_data", "startTime", filter.StartTime.AsTime(), "endTime", filter.EndTime.AsTime(), "eventName", filter.EventName, "userId", filter.UserId, "sessionId", filter.SessionId, "apiKeyId", filter.ApiKeyId)
		nextPageToken = s.pageEncoder.Marshal(filter)
	}

	return &frontendv1.ListAuditLogEventsResponse{
		AuditLogEvents: results,
		NextPageToken:  nextPageToken,
	}, nil
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
