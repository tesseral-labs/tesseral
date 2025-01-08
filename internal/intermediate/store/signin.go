package store

import (
	"context"
	"crypto/sha256"
	"log/slog"
	"time"

	"github.com/google/uuid"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) SignInWithEmail(
	ctx context.Context,
	req *intermediatev1.SignInWithEmailRequest,
) (*intermediatev1.SignInWithEmailResponse, error) {
	projectID := projectid.ProjectID(ctx)

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	token := uuid.New()

	// Send a verification email then issue an intermediate session,
	// so the user can verify their email address and create an organization

	expiresAt := time.Now().Add(15 * time.Minute)
	tokenSha256 := sha256.Sum256(token[:])

	intermediateSession, err := q.CreateIntermediateSession(ctx, queries.CreateIntermediateSessionParams{
		ID:          uuid.Must(uuid.NewV7()),
		ProjectID:   projectID,
		Email:       &req.Email,
		ExpireTime:  &expiresAt,
		TokenSha256: tokenSha256[:],
	})
	if err != nil {
		return nil, err
	}

	// Create a new secret token for the challenge
	secretToken, err := generateSecretToken()
	if err != nil {
		return nil, err
	}
	secretTokenSha256 := sha256.Sum256([]byte(secretToken))

	// TODO: Send the secret token to the user's email address

	evc, err := q.CreateEmailVerificationChallenge(ctx, queries.CreateEmailVerificationChallengeParams{
		ID:                    uuid.New(),
		ChallengeSha256:       secretTokenSha256[:],
		ExpireTime:            &expiresAt,
		IntermediateSessionID: intermediateSession.ID,
		ProjectID:             projectID,
	})
	if err != nil {
		return nil, err
	}

	// TODO: Remove this log line and replace with email sending
	slog.Info("SignInWithEmail", "challenge", secretToken)

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.SignInWithEmailResponse{
		IntermediateSessionToken: idformat.IntermediateSessionToken.Format(token),
		ChallengeId:              idformat.EmailVerificationChallenge.Format(evc.ID),
	}, nil
}
