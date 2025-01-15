package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/errorcodes"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

type EmailVerificationChallenge struct {
	ID                    string
	ChallengeSha256       []byte
	CompleteTime          time.Time
	ProjectID             string
	ExpireTime            time.Time
	IntermediateSessionID string
}

type CreateEmailVerificationChallengeParams struct {
	ChallengeSha256       []byte
	Email                 string
	GoogleUserID          string
	IntermediateSessionID string
	MicrosoftUserID       string
	ProjectID             string
}

type GetEmailVerificationChallengeParams struct {
	Code            string
	Email           string
	GoogleUserID    string
	MicrosoftUserID string
	ProjectID       string
}

func (s *Store) CompleteEmailVerificationChallenge(ctx context.Context, req *intermediatev1.VerifyEmailChallengeRequest) (*intermediatev1.VerifyEmailChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	// Get the email verification challenge from the request
	challengeID, err := idformat.EmailVerificationChallenge.Parse(req.EmailVerificationChallengeId)
	if err != nil {
		return nil, errorcodes.NewInvalidArgumentError(fmt.Errorf("email verification challenge id is invalid"))
	}
	challenge, err := q.GetEmailVerificationChallengeByID(ctx, challengeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorcodes.NewNotFoundError(fmt.Errorf("email verification challenge not found: %w", err))
		}

		return nil, fmt.Errorf("get email verification challenge by id: %w", err)
	}

	// Enforce the intermediate session
	if challenge.IntermediateSessionID.String() != authn.IntermediateSessionID(ctx).String() {
		return nil, errorcodes.NewFailedPreconditionError(fmt.Errorf("intermediate session id mismatch"))
	}

	// Get the intermediate session
	intermediateSession := authn.IntermediateSession(ctx)

	intermediateSessionID, err := idformat.IntermediateSession.Parse(intermediateSession.Id)
	if err != nil {
		return nil, fmt.Errorf("parse intermediate session id: %w", err)
	}

	now := time.Now()
	codeSha256 := sha256.Sum256([]byte(req.Code))

	// Get the email verification challenge
	evc, err := q.GetEmailVerificationChallengeForCompletion(ctx, queries.GetEmailVerificationChallengeForCompletionParams{
		ExpireTime:            &now,
		IntermediateSessionID: intermediateSessionID,
		ProjectID:             projectID,
	})
	if err != nil {
		return nil, err
	}

	err = verifyChallenge(ctx, &evc, codeSha256[:], q)
	if err != nil {
		if err := commit(); err != nil {
			return nil, fmt.Errorf("commit after verify failure: %w", err)
		}

		return nil, errorcodes.NewFailedPreconditionError(fmt.Errorf("verify challenge failure: %w", err))
	}

	// Complete the email verification challenge
	evc, err = q.CompleteEmailVerificationChallenge(ctx, queries.CompleteEmailVerificationChallengeParams{
		CompleteTime: &now,
		ID:           evc.ID,
	})
	if err != nil {
		return nil, err
	}

	// Create a verified email record
	_, err = q.CreateVerifiedEmail(ctx, queries.CreateVerifiedEmailParams{
		ID:                 uuid.New(),
		Email:              intermediateSession.Email,
		GoogleUserID:       &intermediateSession.GoogleUserId,
		GoogleHostedDomain: &intermediateSession.GoogleHostedDomain,
		MicrosoftUserID:    &intermediateSession.MicrosoftUserId,
		MicrosoftTenantID:  &intermediateSession.MicrosoftTenantId,
		ProjectID:          projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("create verified email: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.VerifyEmailChallengeResponse{}, nil
}

func (s *Store) IssueEmailVerificationChallenge(ctx context.Context) (*intermediatev1.IssueEmailVerificationChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	intermediateSessionID := authn.IntermediateSessionID(ctx)

	// Create a new secret token for the challenge
	secretToken, err := generateSecretToken()
	if err != nil {
		return nil, fmt.Errorf("generate secret token: %w", err)
	}
	secretTokenSha256 := sha256.Sum256([]byte(secretToken))

	expiresAt := time.Now().Add(15 * time.Minute)

	// TODO: Send the secret token to the user's email address

	evc, err := q.CreateEmailVerificationChallenge(ctx, queries.CreateEmailVerificationChallengeParams{
		ID:                    uuid.New(),
		ChallengeSha256:       secretTokenSha256[:],
		ExpireTime:            &expiresAt,
		IntermediateSessionID: intermediateSessionID,
		ProjectID:             projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("create email verification challenge: %w", err)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	// TODO: Remove this log line and replace with email sending
	slog.InfoContext(ctx, "IssueEmailVerificationChallenge", "challenge", secretToken)

	return &intermediatev1.IssueEmailVerificationChallengeResponse{
		EmailVerificationChallengeId: idformat.EmailVerificationChallenge.Format(evc.ID),
	}, nil
}

func generateSecretToken() (string, error) {
	// Define the range for a 6-digit number: [100000, 999999]
	min := 100000
	max := 999999

	// Generate a secure random number
	randomNumber := rand.IntN(max-min+1) + min

	return strconv.Itoa(randomNumber), nil
}

func verifyChallenge(ctx context.Context, evc *queries.EmailVerificationChallenge, secretTokenSha256 []byte, q *queries.Queries) error {
	// Check if the challenge has been revoked
	if evc.Revoked {
		return errorcodes.NewFailedPreconditionError(fmt.Errorf("email verification challenge revoked"))
	}

	// Check if the challenge has already been completed
	if evc.CompleteTime != nil {
		_, err := q.RevokeEmailVerificationChallenge(ctx, evc.ID)
		if err != nil {
			return fmt.Errorf("revoke email verification challenge after complete time failure: %w", err)
		}

		return errorcodes.NewFailedPreconditionError(fmt.Errorf("email verification challenge already completed"))
	}

	// Check if the challenge has expired
	if evc.ExpireTime.Before(time.Now()) {
		_, err := q.RevokeEmailVerificationChallenge(ctx, evc.ID)
		if err != nil {
			return fmt.Errorf("revoke email verification challenge after expire time failure: %w", err)
		}

		return errorcodes.NewFailedPreconditionError(fmt.Errorf("email verification challenge expired"))
	}

	// Check if the challenge is correct
	if !bytes.Equal(evc.ChallengeSha256, secretTokenSha256) {
		_, err := q.RevokeEmailVerificationChallenge(ctx, evc.ID)
		if err != nil {
			return fmt.Errorf("revoke email verification challenge: %w", err)
		}

		return errorcodes.NewFailedPreconditionError(fmt.Errorf("email verification challenge code mismatch"))
	}

	return nil
}
