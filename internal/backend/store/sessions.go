package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/shared/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) ListSessions(ctx context.Context, req *backendv1.ListSessionsRequest) (*backendv1.ListSessionsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.UserId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	// authz
	if _, err := q.GetUser(ctx, queries.GetUserParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        userID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qSessions, err := q.ListSessions(ctx, queries.ListSessionsParams{
		UserID: userID,
		ID:     startID,
		Limit:  int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	var sessions []*backendv1.Session
	for _, qSession := range qSessions {
		sessions = append(sessions, parseSession(qSession))
	}

	var nextPageToken string
	if len(sessions) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qSessions[limit].ID)
		sessions = sessions[:limit]
	}

	return &backendv1.ListSessionsResponse{
		Sessions:      sessions,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetSession(ctx context.Context, req *backendv1.GetSessionRequest) (*backendv1.GetSessionResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	sessionID, err := idformat.Session.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid session id", fmt.Errorf("parse session id: %w", err))
	}

	qSession, err := q.GetSession(ctx, queries.GetSessionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        sessionID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("session not found", fmt.Errorf("get session: %w", err))
		}

		return nil, fmt.Errorf("get session: %w", err)
	}

	return &backendv1.GetSessionResponse{Session: parseSession(qSession)}, nil
}

func parseSession(qSession queries.Session) *backendv1.Session {
	return &backendv1.Session{
		Id:      idformat.Session.Format(qSession.ID),
		UserId:  idformat.User.Format(qSession.UserID),
		Revoked: qSession.Revoked,
	}
}
