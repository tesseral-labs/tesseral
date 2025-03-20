package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) ExchangeRelayedSessionTokenForSession(ctx context.Context, req *intermediatev1.ExchangeRelayedSessionTokenForSessionRequest) (*intermediatev1.ExchangeRelayedSessionTokenForSessionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	relayedSessionToken, err := idformat.RelayedSessionToken.Parse(req.RelayedSessionToken)
	if err != nil {
		return nil, fmt.Errorf("parse relayed session token: %w", err)
	}

	relayedSessionTokenSHA := sha256.Sum256(relayedSessionToken[:])
	qRelayedSession, err := q.GetRelayedSessionByTokenSHA256(ctx, relayedSessionTokenSHA[:])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewUnauthenticatedError("invalid relayed session token", fmt.Errorf("invalid relayed session token"))
		}

		return nil, fmt.Errorf("get relayed session by token sha256: %w", err)
	}

	relayedRefreshTokenUUID := uuid.New()
	relayedRefreshTokenSHA := sha256.Sum256(relayedRefreshTokenUUID[:])
	if _, err := q.UpdateRelayedSessionRefreshTokenSHA256(ctx, queries.UpdateRelayedSessionRefreshTokenSHA256Params{
		SessionID:                 qRelayedSession.SessionID,
		RelayedRefreshTokenSha256: relayedRefreshTokenSHA[:],
	}); err != nil {
		return nil, fmt.Errorf("update relayed session refresh token sha256: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.ExchangeRelayedSessionTokenForSessionResponse{
		RefreshToken:        idformat.RelayedSessionRefreshToken.Format(relayedRefreshTokenUUID),
		RelayedSessionState: derefOrEmpty(qRelayedSession.State),
		AccessToken:         "", // populated in service
	}, nil
}
