package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

var emailVerificationChallengeDuration = time.Minute * 10

func (s *Store) IssueEmailVerificationChallenge(ctx context.Context, req *intermediatev1.IssueEmailVerificationChallengeRequest) (*intermediatev1.IssueEmailVerificationChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	if qIntermediateSession.Email != nil && *qIntermediateSession.Email != req.Email {
		return nil, apierror.NewInvalidArgumentError("email does not match existing value on intermediate session", fmt.Errorf("email does not match existing value on intermediate session"))
	}

	if qIntermediateSession.Email == nil {
		if _, err := q.UpdateIntermediateSessionEmail(ctx, queries.UpdateIntermediateSessionEmailParams{
			ID:    authn.IntermediateSessionID(ctx),
			Email: &req.Email,
		}); err != nil {
			return nil, fmt.Errorf("update intermediate session email: %w", err)
		}
	}

	secretToken := generateSecretToken()
	secretTokenSHA256 := sha256.Sum256([]byte(secretToken))

	expireTime := time.Now().Add(emailVerificationChallengeDuration)

	qEmailVerificationChallenge, err := q.CreateEmailVerificationChallenge(ctx, queries.CreateEmailVerificationChallengeParams{
		ID:                    uuid.New(),
		ChallengeSha256:       secretTokenSHA256[:],
		ExpireTime:            &expireTime,
		IntermediateSessionID: authn.IntermediateSessionID(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("create email verification challenge: %w", err)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	// TODO: Remove this log line and replace with email sending
	slog.InfoContext(ctx, "TODO-REMOVEME-IssueEmailVerificationChallenge", "challenge", secretToken)

	return &intermediatev1.IssueEmailVerificationChallengeResponse{
		EmailVerificationChallengeId: idformat.EmailVerificationChallenge.Format(qEmailVerificationChallenge.ID),
	}, nil
}

func (s *Store) VerifyEmailChallenge(ctx context.Context, req *intermediatev1.VerifyEmailChallengeRequest) (*intermediatev1.VerifyEmailChallengeResponse, error) {
	intermediateSession := authn.IntermediateSession(ctx)

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	challengeSHA256 := sha256.Sum256([]byte(req.Code))
	qEmailVerificationChallenge, err := q.GetEmailVerificationChallengeByChallengeSHA(ctx, queries.GetEmailVerificationChallengeByChallengeSHAParams{
		IntermediateSessionID: authn.IntermediateSessionID(ctx),
		ChallengeSha256:       challengeSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("get email verification challenge by challenge sha: %w", err)
	}

	if _, err := q.CompleteEmailVerificationChallenge(ctx, qEmailVerificationChallenge.ID); err != nil {
		return nil, fmt.Errorf("complete email verification challenge: %w", err)
	}

	if intermediateSession.GoogleUserId != "" || intermediateSession.MicrosoftUserId != "" {
		if _, err := q.CreateVerifiedEmail(ctx, queries.CreateVerifiedEmailParams{
			ID:              uuid.New(),
			ProjectID:       authn.ProjectID(ctx),
			Email:           authn.IntermediateSession(ctx).Email,
			GoogleUserID:    refOrNil(intermediateSession.GoogleUserId),
			MicrosoftUserID: refOrNil(intermediateSession.MicrosoftUserId),
		}); err != nil {
			return nil, fmt.Errorf("create verified email: %w", err)
		}
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.VerifyEmailChallengeResponse{}, nil
}

func generateSecretToken() string {
	// Define the range for a 6-digit number: [100000, 999999]
	min := 100000
	max := 999999

	// Generate a secure random number
	randomNumber := rand.IntN(max-min+1) + min

	return strconv.Itoa(randomNumber)
}
