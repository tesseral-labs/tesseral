package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"html/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
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

	emailVerificationChallengeCode := uuid.New()
	secretTokenSHA256 := sha256.Sum256(emailVerificationChallengeCode[:])

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

	if err := s.sendEmailVerificationChallenge(ctx, req.Email, idformat.EmailVerificationChallengeCode.Format(emailVerificationChallengeCode)); err != nil {
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

	emailVerificationChallengeCodeUUID, err := idformat.EmailVerificationChallengeCode.Parse(req.Code)
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

var emailVerificationEmailBodyTmpl = template.Must(template.New("emailVerificationEmailBody").Parse(`
Hi,

To continue logging in to {{ .ProjectDisplayName }}, please verify your email address by visiting the link below.

{{ .EmailVerificationLink }}

If you did not request this verification, please ignore this email.
`))

func (s *Store) sendEmailVerificationChallenge(ctx context.Context, toAddress string, secretToken string) error {
	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	subject := fmt.Sprintf("%s - Verify your email address", qProject.DisplayName)

	var body bytes.Buffer
	emailVerificationEmailBodyTmpl.Execute(&body, struct {
		ProjectDisplayName    string
		EmailVerificationLink string
	}{
		ProjectDisplayName:    qProject.DisplayName,
		EmailVerificationLink: fmt.Sprintf("https://%s/verify-email?code=%s", qProject.VaultDomain, secretToken),
	})

	if _, err := s.ses.SendEmail(ctx, &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: &subject,
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(body.String()),
					},
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{toAddress},
		},
		FromEmailAddress: aws.String(fmt.Sprintf("noreply@%s", qProject.EmailSendFromDomain)),
	}); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
