package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"log/slog"
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

var ErrProjectIDRequired = errors.New("project ID required")

type EmailVerificationChallenge struct {
	ID              string
	ProjectID       string
	Challenge       string
	ChallengeSha256 []byte
	CompleteTime    time.Time
	Email           string
	ExpireTime      time.Time
	GoogleUserID    string
	MicrosoftUserID string
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

	intermediateSessionCtx := authn.IntermediateSession(ctx)
	intermediateSessionID, err := idformat.IntermediateSession.Parse(intermediateSessionCtx.Id)
	if err != nil {
		return nil, err
	}

	intermediateSession, err := q.GetIntermediateSessionByID(ctx, intermediateSessionID)
	if err != nil {
		return nil, err
	}

	slog.Info("CompleteEmailVerificationChallenge", "intermediateSession", intermediateSessionID)

	now := time.Now()
	secretTokenSha256 := sha256.Sum256([]byte(challenge))

	params := queries.GetEmailVerificationChallengeParams{
		ExpireTime:      &now,
		ChallengeSha256: secretTokenSha256[:],
		Email:           intermediateSession.Email,
		GoogleUserID:    intermediateSession.GoogleUserID,
		MicrosoftUserID: intermediateSession.MicrosoftUserID,
		ProjectID:       projectID,
	}

	slog.Info("CompleteEmailVerificationChallenge", "params", params)

	evc, err := q.GetEmailVerificationChallenge(ctx, params)
	if err != nil {
		return nil, err
	}

	evc, err = q.CompleteEmailVerificationChallenge(ctx, queries.CompleteEmailVerificationChallengeParams{
		CompleteTime: &now,
		ID:           evc.ID,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.VerifyEmailChallengeResponse{
		Success: true,
	}, nil
}

func (s *Store) CreateEmailVerificationChallenge(ctx context.Context, params *CreateEmailVerificationChallengeParams) (*EmailVerificationChallenge, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(params.ProjectID)
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

	expiresAt := time.Now().Add(15 * time.Minute)

	evc, err := q.CreateEmailVerificationChallenge(ctx, queries.CreateEmailVerificationChallengeParams{
		ID:              uuid.New(),
		ProjectID:       projectID,
		ChallengeSha256: secretTokenSha256[:],
		Email:           &params.Email,
		ExpireTime:      &expiresAt,
		GoogleUserID:    &params.GoogleUserID,
		MicrosoftUserID: &params.MicrosoftUserID,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseEmailVerificationChallenge(&evc, secretToken), nil
}

func (s *Store) GetEmailVerificationChallenge(ctx context.Context, params *GetEmailVerificationChallengeParams) (*EmailVerificationChallenge, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	now := time.Now()

	projectID, err := idformat.Project.Parse(params.ProjectID)
	if err != nil {
		return nil, err
	}

	secretTokenSha256 := sha256.Sum256([]byte(params.Code))

	evc, err := q.GetEmailVerificationChallenge(ctx, queries.GetEmailVerificationChallengeParams{
		ExpireTime:      &now,
		ChallengeSha256: secretTokenSha256[:],
		Email:           &params.Email,
		GoogleUserID:    &params.GoogleUserID,
		MicrosoftUserID: &params.MicrosoftUserID,
		ProjectID:       projectID,
	})
	if err != nil {
		return nil, err
	}

	return parseEmailVerificationChallenge(&evc, ""), nil
}

func generateSecretToken() (string, error) {
	// Define the range for a 6-digit number: [100000, 999999]
	min := 100000
	max := 999999

	// Generate a secure random number
	randomNumber := rand.IntN(max-min+1) + min

	return strconv.Itoa(randomNumber), nil
}

func parseEmailVerificationChallenge(evc *queries.EmailVerificationChallenge, originalChallenge string) *EmailVerificationChallenge {
	return &EmailVerificationChallenge{
		ID:              idformat.EmailVerificationChallenge.Format(evc.ID),
		ProjectID:       idformat.Project.Format(evc.ProjectID),
		Challenge:       originalChallenge,
		ChallengeSha256: evc.ChallengeSha256,
		CompleteTime:    *evc.CompleteTime,
		Email:           *evc.Email,
		ExpireTime:      *evc.ExpireTime,
		GoogleUserID:    *evc.GoogleUserID,
		MicrosoftUserID: *evc.MicrosoftUserID,
	}
}
