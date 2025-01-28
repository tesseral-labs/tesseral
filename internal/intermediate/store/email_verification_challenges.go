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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
)

func (s *Store) IssueEmailVerificationChallenge(ctx context.Context, req *intermediatev1.IssueEmailVerificationChallengeRequest) (*intermediatev1.IssueEmailVerificationChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("get project by id: %w", fmt.Errorf("project not found: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

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

	_, err = q.UpdateIntermediateSessionEmailVerificationChallengeSha256(ctx, queries.UpdateIntermediateSessionEmailVerificationChallengeSha256Params{
		ID:                               authn.IntermediateSessionID(ctx),
		EmailVerificationChallengeSha256: secretTokenSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("set email verification challenge: %w", err)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	err = s.sendEmailVerificationChallenge(ctx, authn.IntermediateSession(ctx).Email, secretToken)
	if err != nil {
		return nil, fmt.Errorf("send email verification challenge: %w", err)
	}

	return &intermediatev1.IssueEmailVerificationChallengeResponse{}, nil
}

func (s *Store) VerifyEmailChallenge(ctx context.Context, req *intermediatev1.VerifyEmailChallengeRequest) (*intermediatev1.VerifyEmailChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("get intermediate session by id: %w", fmt.Errorf("intermediate session not found: %w", err))
		}

		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("get project by id: %w", fmt.Errorf("project not found: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	challengeSHA256 := sha256.Sum256([]byte(req.Code))
	if !bytes.Equal(qIntermediateSession.EmailVerificationChallengeSha256, challengeSHA256[:]) {
		return nil, apierror.NewInvalidArgumentError("invalid email verification code", fmt.Errorf("invalid email verification code"))
	}

	if _, err := q.UpdateIntermediateSessionEmailVerificationChallengeCompleted(ctx, authn.IntermediateSessionID(ctx)); err != nil {
		return nil, fmt.Errorf("update intermediate session email verified: %w", err)
	}

	if qIntermediateSession.GoogleUserID != nil || qIntermediateSession.MicrosoftUserID != nil {
		if _, err := q.CreateVerifiedEmail(ctx, queries.CreateVerifiedEmailParams{
			ID:              uuid.New(),
			ProjectID:       authn.ProjectID(ctx),
			Email:           *qIntermediateSession.Email,
			GoogleUserID:    qIntermediateSession.GoogleUserID,
			MicrosoftUserID: qIntermediateSession.MicrosoftUserID,
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

func (s *Store) sendEmailVerificationChallenge(ctx context.Context, email string, secretToken string) error {
	output, err := s.ses.SendEmail(ctx, &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(fmt.Sprintf("<h2>Please verifiy your email address to continue logging in</h2><p>Your email verification code is: %s</p>", secretToken)),
					},
				},
				Subject: &types.Content{
					Data: aws.String("Verify your email address"),
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		FromEmailAddress: aws.String("replace-me@tesseral.app"), // TODO: Replace with a real email address once verification is in place
	})
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	slog.InfoContext(ctx, "sendEmailVerificationChallenge", "output", output)

	return nil
}
