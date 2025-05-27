package store

import (
	"context"
	"fmt"

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
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	var filter *frontendv1.ListAuditLogEventsRequest_Filter
	if req.PageToken != "" {
		filter = new(frontendv1.ListAuditLogEventsRequest_Filter)
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
	startID := uuidv7.NilWithTime(startTime)

	wheres := []sq.Sqlizer{
		sq.Eq{"project_id": projectID[:]},
		sq.Eq{"organization_id": orgID[:]},
		sq.Gt{"id": startID},
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
	if len(results) > 0 {
		last := results[len(results)-1]
		filter.StartTime = last.EventTime
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
