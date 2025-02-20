package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) RedeemUserImpersonationToken(ctx context.Context, req *intermediatev1.RedeemUserImpersonationTokenRequest) (*intermediatev1.RedeemUserImpersonationTokenResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	secretTokenUUID, err := idformat.UserImpersonationSecretToken.Parse(req.SecretUserImpersonationToken)
	if err != nil {
		return nil, fmt.Errorf("parse user impersonation secret token: %w", err)
	}

	secretTokenSHA256 := sha256.Sum256(secretTokenUUID[:])
	qUserImpersonationToken, err := q.GetUserImpersonationTokenBySecretTokenSHA256(ctx, secretTokenSHA256[:])
	if err != nil {
		return nil, fmt.Errorf("get user impersonation token by secret token sha256: %w", err)
	}

	expireTime := time.Now().Add(sessionDuration)

	// Create a new session for the user
	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])
	if _, err := q.CreateImpersonatedSession(ctx, queries.CreateImpersonatedSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expireTime,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qUserImpersonationToken.ImpersonatedID,
		ImpersonatorUserID: &qUserImpersonationToken.ImpersonatorID,
	}); err != nil {
		return nil, fmt.Errorf("create impersonated session: %w", err)
	}

	if _, err := q.RevokeUserImpersonationToken(ctx, qUserImpersonationToken.ID); err != nil {
		return nil, fmt.Errorf("revoke user impersonation token: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RedeemUserImpersonationTokenResponse{
		AccessToken:  "", // populated in service
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}
