package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type OpenAuthSession struct {
	ID 					uuid.UUID
	UserID 			uuid.UUID
	CreateTime 	time.Time
	ExpireTime 	time.Time
	Token 			string
	TokenSha256 []byte
	Revoked 		bool
}

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionExpired = errors.New("session expired")

type CreateSessionRequest struct {
	UserID string
}

func (s *Store) CreateSession(ctx context.Context, req *CreateSessionRequest) (*OpenAuthSession, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return &OpenAuthSession{}, err
	}
	defer rollback()

	userId, err := idformat.User.Parse(req.UserID)
	if err != nil {
		return &OpenAuthSession{}, err
	}

	// Sessions expire after 7 days
	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	session, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID: uuid.New(),
		UserID: userId,
		ExpireTime: &expiresAt,
		Token: uuid.New().String(),
	})
	if err != nil {
		return &OpenAuthSession{}, err
	}

	if err := commit(); err != nil {
		return &OpenAuthSession{}, err
	}

	return transformSession(session), nil
}

func (s *Store) RevokeSession() error {
	return nil
}

func transformSession(session queries.Session) *OpenAuthSession {
	return &OpenAuthSession{
		ID: session.ID,
		UserID: session.UserID,
		CreateTime: *session.CreateTime,
		ExpireTime: *session.ExpireTime,
		Token: session.Token,
		Revoked: session.Revoked,
	}
}