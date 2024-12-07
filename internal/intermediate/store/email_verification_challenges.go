package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

var ErrEmailVerificationChallengeMismatch = errors.New("email verification challenge mismatch")
var ErrEmailVerficationChallengeNotFound = errors.New("email verification challenge not found")
var ErrEmailVerificationChallengeInvalidState = errors.New("email verification challenge in invalid state")
var ErrEmailVerificationChallengeExpired = errors.New("email verification challenge expired")
var ErrIntermediateSessionIDRequired = errors.New("intermediate session ID required")
var ErrIntermediateSessionRequired = errors.New("intermediate session required")
var ErrProjectIDRequired = errors.New("project ID required")

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

func (s *Store) CompleteEmailVerificationChallenge(ctx context.Context, challenge string) (*intermediatev1.VerifyEmailChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := projectid.ProjectID(ctx)
	if projectID == uuid.Nil {
		return nil, ErrProjectIDRequired
	}

	// Get the email verification challenge ID from the intermediate session
	intermediateSession := authn.IntermediateSession(ctx)
	if intermediateSession == nil {
		return nil, ErrIntermediateSessionRequired
	}

	intermediateSessionID, err := idformat.IntermediateSession.Parse(intermediateSession.Id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	secretTokenSha256 := sha256.Sum256([]byte(challenge))

	// Get the email verification challenge
	evc, err := q.GetEmailVerificationChallengeForCompletion(ctx, queries.GetEmailVerificationChallengeForCompletionParams{
		ExpireTime:            &now,
		IntermediateSessionID: intermediateSessionID,
		ProjectID:             projectID,
	})
	if err != nil {
		return nil, err
	}

	err = mustVerifyChallenge(ctx, &evc, secretTokenSha256[:], q)
	if err != nil {
		if err := commit(); err != nil {
			return nil, err
		}

		return nil, err
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
		GoogleUserID:       intermediateSession.GoogleUserId,
		GoogleHostedDomain: intermediateSession.GoogleHostedDomain,
		MicrosoftUserID:    intermediateSession.MicrosoftUserId,
		MicrosoftTenantID:  intermediateSession.MicrosoftTenantId,
		ProjectID:          projectID,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.VerifyEmailChallengeResponse{}, nil
}

func (s *Store) IssueEmailVerificationChallenge(ctx context.Context) (*intermediatev1.IssueEmailVerificationChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := projectid.ProjectID(ctx)
	if projectID == uuid.Nil {
		return nil, ErrProjectIDRequired
	}

	intermediateSessionID := authn.IntermediateSessionID(ctx)
	if intermediateSessionID == uuid.Nil {
		return nil, ErrIntermediateSessionIDRequired
	}

	// Create a new secret token for the challenge
	secretToken, err := generateSecretToken()
	if err != nil {
		return nil, err
	}
	secretTokenSha256 := sha256.Sum256([]byte(secretToken))

	expiresAt := time.Now().Add(15 * time.Minute)

	// TODO: Send the secret token to the user's email address

	_, err = q.CreateEmailVerificationChallenge(ctx, queries.CreateEmailVerificationChallengeParams{
		ID:                    uuid.New(),
		ChallengeSha256:       secretTokenSha256[:],
		ExpireTime:            &expiresAt,
		IntermediateSessionID: intermediateSessionID,
		ProjectID:             projectID,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.IssueEmailVerificationChallengeResponse{}, nil
}

func generateSecretToken() (string, error) {
	// Define the range for a 6-digit number: [100000, 999999]
	min := 100000
	max := 999999

	// Generate a secure random number
	randomNumber := rand.IntN(max-min+1) + min

	return strconv.Itoa(randomNumber), nil
}

func mustVerifyChallenge(ctx context.Context, evc *queries.EmailVerificationChallenge, secretTokenSha256 []byte, q *queries.Queries) error {
	// Check if the challenge has been revoked
	if evc.Revoked {
		return ErrEmailVerificationChallengeInvalidState
	}

	// Check if the challenge has already been completed
	if evc.CompleteTime != nil {
		err := revokeEmailVerificationChallenge(ctx, evc.ID, q)
		if err != nil {
			return err
		}

		return ErrEmailVerificationChallengeInvalidState
	}

	// Check if the challenge has expired
	if evc.ExpireTime.Before(time.Now()) {
		err := revokeEmailVerificationChallenge(ctx, evc.ID, q)
		if err != nil {
			return err
		}

		return ErrEmailVerificationChallengeInvalidState
	}

	// Check if the challenge is correct
	if !bytes.Equal(evc.ChallengeSha256, secretTokenSha256) {
		err := revokeEmailVerificationChallenge(ctx, evc.ID, q)
		if err != nil {
			return err
		}

		return ErrEmailVerificationChallengeMismatch
	}

	return nil
}

func parseEmailVerificationChallenge(evc *queries.EmailVerificationChallenge, originalChallenge string) *EmailVerificationChallenge {
	return &EmailVerificationChallenge{
		ID:                    idformat.EmailVerificationChallenge.Format(evc.ID),
		ChallengeSha256:       evc.ChallengeSha256,
		CompleteTime:          *evc.CompleteTime,
		IntermediateSessionID: idformat.IntermediateSession.Format(evc.IntermediateSessionID),
		ProjectID:             idformat.Project.Format(evc.ProjectID),
	}
}

func revokeEmailVerificationChallenge(ctx context.Context, id uuid.UUID, q *queries.Queries) error {
	_, err := q.RevokeEmailVerificationChallenge(ctx, id)

	return err
}
