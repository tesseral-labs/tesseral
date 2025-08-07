package emailworker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riverqueue/river"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store"
)

type Worker struct {
	Store *store.Store
	river.WorkerDefaults[Args]
}

type Args struct {
	ProjectID     string
	VerifyEmail   *VerifyEmailParams
	PasswordReset *PasswordResetParams
	UserInvite    *UserInviteParams
}

func (Args) Kind() string {
	return "email"
}

type VerifyEmailParams struct {
	EmailAddress          string
	EmailVerificationCode string
}

type PasswordResetParams struct {
	EmailAddress      string
	PasswordResetCode string
}

type UserInviteParams struct {
	UserInviteID string
}

func (w *Worker) Work(ctx context.Context, job *river.Job[Args]) error {
	slog.InfoContext(ctx, "work", "project_id", job.Args.ProjectID)

	switch {
	case job.Args.VerifyEmail != nil:
		if err := w.Store.SendEmailVerifyEmail(ctx, &store.SendEmailVerifyEmailRequest{
			ProjectID:             job.Args.ProjectID,
			EmailAddress:          job.Args.VerifyEmail.EmailAddress,
			EmailVerificationCode: job.Args.VerifyEmail.EmailVerificationCode,
		}); err != nil {
			return fmt.Errorf("send verify email: %w", err)
		}
	case job.Args.PasswordReset != nil:
		if err := w.Store.SendEmailPasswordReset(ctx, &store.SendEmailPasswordResetRequest{
			ProjectID:         job.Args.ProjectID,
			EmailAddress:      job.Args.PasswordReset.EmailAddress,
			PasswordResetCode: job.Args.PasswordReset.PasswordResetCode,
		}); err != nil {
			return fmt.Errorf("send password reset: %w", err)
		}
	case job.Args.UserInvite != nil:
		if err := w.Store.SendEmailUserInvite(ctx, &store.SendEmailUserInviteRequest{
			ProjectID:    job.Args.ProjectID,
			UserInviteID: job.Args.UserInvite.UserInviteID,
		}); err != nil {
			return fmt.Errorf("send user invite: %w", err)
		}
	}

	return nil
}
