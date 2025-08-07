package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/emailworker"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

// defaultEmailQuotaDaily is the default number of emails a project may send per
// day.
var defaultEmailQuotaDaily int32 = 1000

func (s *Store) IssueEmailVerificationChallenge(ctx context.Context, req *intermediatev1.IssueEmailVerificationChallengeRequest) (*intermediatev1.IssueEmailVerificationChallengeResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
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

	emailVerificationChallengeCode := uuid.New()
	secretTokenSHA256 := sha256.Sum256(emailVerificationChallengeCode[:])

	_, err = q.UpdateIntermediateSessionEmailVerificationChallengeSha256(ctx, queries.UpdateIntermediateSessionEmailVerificationChallengeSha256Params{
		ID:                               authn.IntermediateSessionID(ctx),
		EmailVerificationChallengeSha256: secretTokenSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("set email verification challenge: %w", err)
	}

	qEmailDailyQuotaUsage, err := q.IncrementProjectEmailDailyQuotaUsage(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("increment project email daily quota usage: %w", err)
	}

	emailQuotaDaily := defaultEmailQuotaDaily
	if qProject.EmailQuotaDaily != nil {
		emailQuotaDaily = *qProject.EmailQuotaDaily
	}

	slog.InfoContext(ctx, "email_daily_quota_usage", "usage", qEmailDailyQuotaUsage.QuotaUsage, "quota", emailQuotaDaily)

	if qEmailDailyQuotaUsage.QuotaUsage > emailQuotaDaily {
		slog.InfoContext(ctx, "email_daily_quota_exceeded")
		return nil, apierror.NewFailedPreconditionError("email daily quota exceeded", fmt.Errorf("email daily quota exceeded"))
	}

	jobInsertRes, err := s.riverClient.InsertTx(ctx, tx, emailworker.Args{
		ProjectID: idformat.Project.Format(authn.ProjectID(ctx)),
		VerifyEmail: &emailworker.VerifyEmailParams{
			EmailAddress:          req.Email,
			EmailVerificationCode: idformat.EmailVerificationChallengeCode.Format(emailVerificationChallengeCode),
		},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("insert email worker job: %w", err)
	}

	slog.InfoContext(ctx, "email_worker_job_inserted", "job_id", jobInsertRes.Job.ID)

	if err := commit(); err != nil {
		return nil, err
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

	emailVerificationChallengeCodeUUID, err := idformat.EmailVerificationChallengeCode.Parse(strings.TrimSpace(req.Code))
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid email verification code", fmt.Errorf("invalid email verification code"))
	}

	challengeSHA256 := sha256.Sum256(emailVerificationChallengeCodeUUID[:])
	if !bytes.Equal(qIntermediateSession.EmailVerificationChallengeSha256, challengeSHA256[:]) {
		return nil, apierror.NewInvalidArgumentError("invalid email verification code", fmt.Errorf("invalid email verification code"))
	}

	if _, err := q.UpdateIntermediateSessionEmailVerificationChallengeCompleted(ctx, authn.IntermediateSessionID(ctx)); err != nil {
		return nil, fmt.Errorf("update intermediate session email verified: %w", err)
	}

	if qIntermediateSession.GoogleUserID != nil || qIntermediateSession.MicrosoftUserID != nil || qIntermediateSession.GithubUserID != nil {
		if _, err := q.CreateVerifiedEmail(ctx, queries.CreateVerifiedEmailParams{
			ID:              uuid.New(),
			ProjectID:       authn.ProjectID(ctx),
			Email:           *qIntermediateSession.Email,
			GoogleUserID:    qIntermediateSession.GoogleUserID,
			MicrosoftUserID: qIntermediateSession.MicrosoftUserID,
			GithubUserID:    qIntermediateSession.GithubUserID,
		}); err != nil {
			return nil, fmt.Errorf("create verified email: %w", err)
		}
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.VerifyEmailChallengeResponse{}, nil
}
