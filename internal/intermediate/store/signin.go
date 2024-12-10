package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) SignInWithEmail(
	ctx *context.Context,
	req *intermediatev1.SignInWithEmailRequest,
) (*intermediatev1.SignInWithEmailResponse, error) {
	projectID := projectid.ProjectID(*ctx)

	shouldVerify, err := s.shouldVerifyEmail(*ctx, projectID, req.Email, "", "")
	if err != nil {
		return nil, err
	}

	_, q, commit, rollback, err := s.tx(*ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	users, err := q.ListUsersByEmail(*ctx, req.Email)
	if err != nil {
		return nil, err
	}

	token := uuid.New()

	if users != nil {
		// TODO: Implement factor checking before issuing a session
		panic(errors.New("not implemented"))
	} else {
		// Send a verification email then issue an intermediate session,
		// so the user can verify their email address and create an organization

		expiresAt := time.Now().Add(15 * time.Minute)
		tokenSha256 := sha256.Sum256(token[:])

		intermediateSession, err := q.CreateIntermediateSession(*ctx, queries.CreateIntermediateSessionParams{
			ID:          uuid.New(),
			ProjectID:   projectID,
			Email:       &req.Email,
			ExpireTime:  &expiresAt,
			TokenSha256: tokenSha256[:],
		})
		if err != nil {
			return nil, err
		}

		// TODO: Remove this after we're handling cookies properly
		slog.Info("SignInWithEmail", "intermediate_session_token", idformat.IntermediateSessionToken.Format(token))

		var evcid *string = nil

		if shouldVerify {
			// Create a new secret token for the challenge
			secretToken, err := generateSecretToken()
			if err != nil {
				return nil, err
			}
			secretTokenSha256 := sha256.Sum256([]byte(secretToken))

			// TODO: Send the secret token to the user's email address

			expiresAt := time.Now().Add(15 * time.Minute)

			evc, err := q.CreateEmailVerificationChallenge(*ctx, queries.CreateEmailVerificationChallengeParams{
				ID:                    uuid.New(),
				ChallengeSha256:       secretTokenSha256[:],
				ExpireTime:            &expiresAt,
				IntermediateSessionID: intermediateSession.ID,
				ProjectID:             projectID,
			})
			if err != nil {
				return nil, err
			}

			evcID := idformat.EmailVerificationChallenge.Format(evc.ID)
			evcid = &evcID

			// TODO: Remove this log line and replace with email sending
			slog.Info("SignInWithEmail", "challenge", secretToken)
		}

		if err := commit(); err != nil {
			return nil, err
		}

		return &intermediatev1.SignInWithEmailResponse{
			ChallengeID: *evcid,
		}, nil
	}
}

func (s *Store) shouldVerifyEmail(ctx context.Context, projectID uuid.UUID, email string, googleUserID string, microsoftUserID string) (bool, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return false, err
	}
	defer rollback()

	verifiedEmails, err := q.ListVerifiedEmails(ctx, queries.ListVerifiedEmailsParams{
		ProjectID:       projectID,
		Email:           email,
		GoogleUserID:    &googleUserID,
		MicrosoftUserID: &microsoftUserID,
	})
	if err != nil {
		return false, err
	}

	if len(verifiedEmails) == 0 {
		return true, nil
	} else {
		return false, nil
	}
}
