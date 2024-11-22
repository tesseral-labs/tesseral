package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type Session struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CreateTime time.Time
	ExpireTime time.Time
	Revoked    bool
}

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionExpired = errors.New("session expired")

type CreateSessionRequest struct {
	UserID string
}

func (s *Store) CreateSession(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return &Session{}, err
	}
	defer rollback()

	userId, err := idformat.User.Parse(req.UserID)
	if err != nil {
		return &Session{}, err
	}

	// Sessions expire after 7 days
	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	session, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:         uuid.New(),
		UserID:     userId,
		ExpireTime: &expiresAt,
	})
	if err != nil {
		return &Session{}, err
	}

	if err := commit(); err != nil {
		return &Session{}, err
	}

	return parseSession(session), nil
}

func (s *Store) RevokeSession() error {
	return nil
}

func parseSession(session queries.Session) *Session {
	return &Session{
		ID:         session.ID,
		UserID:     session.UserID,
		CreateTime: *session.CreateTime,
		ExpireTime: *session.ExpireTime,
		Revoked:    session.Revoked,
	}
}
