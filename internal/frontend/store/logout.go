package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Store) Logout(ctx context.Context, req *frontendv1.LogoutRequest) (*frontendv1.LogoutResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}
	defer rollback()

	sessionID := authn.SessionID(ctx)
	qSession, err := q.GetSessionByID(ctx, sessionID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get session by id: %w", err)
	}

	// Invalidate the session if one exists
	var emptyUUID uuid.UUID

	if qSession.ID != emptyUUID {
		if err := q.InvalidateSession(ctx, qSession.ID); err != nil {
			return nil, fmt.Errorf("delete session: %w", err)
		}
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.LogoutResponse{}, nil
}
