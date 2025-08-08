package store

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type SendEmailVerifyEmailRequest struct {
	ProjectID             string
	EmailAddress          string
	EmailVerificationCode string
}

func (s *Store) SendEmailVerifyEmail(ctx context.Context, req *SendEmailVerifyEmailRequest) error {
	projectID, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	qProject, err := s.q().GetProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	if qProject.CustomEmailVerifyEmail {
		if err := s.SendWebhook(ctx, &SendWebhookRequest{
			ProjectID: req.ProjectID,
			EventType: "custom_email.verify_email",
			Payload: map[string]any{
				"type":                           "custom_email.verify_email",
				"emailAddress":                   req.EmailAddress,
				"emailVerificationChallengeCode": req.EmailVerificationCode,
			},
		}); err != nil {
			return fmt.Errorf("send webhook: %w", err)
		}
		return nil
	}

	vaultDomain := qProject.VaultDomain
	if req.ProjectID == s.ConsoleProjectID {
		vaultDomain = s.ConsoleDomain
	}

	var body bytes.Buffer
	if err := emailVerificationEmailBodyTmpl.Execute(&body, struct {
		ProjectDisplayName    string
		EmailVerificationLink string
		EmailVerificationCode string
	}{
		ProjectDisplayName:    qProject.DisplayName,
		EmailVerificationLink: fmt.Sprintf("https://%s/verify-email?code=%s", vaultDomain, req.EmailVerificationCode),
		EmailVerificationCode: req.EmailVerificationCode,
	}); err != nil {
		return fmt.Errorf("execute email verification email body template: %w", err)
	}

	if _, err := s.SES.SendEmail(ctx, &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(fmt.Sprintf("%s - Verify your email address", qProject.DisplayName)),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(body.String()),
					},
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{req.EmailAddress},
		},
		FromEmailAddress: aws.String(fmt.Sprintf("noreply@%s", qProject.EmailSendFromDomain)),
	}); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

var emailVerificationEmailBodyTmpl = template.Must(template.New("emailVerificationEmailBody").Parse(`Hello,

To continue logging in to {{ .ProjectDisplayName }}, please verify your email address by visiting the link below.

{{ .EmailVerificationLink }}

You can also go back to the "Check your email" page and enter this verification code manually:

{{ .EmailVerificationCode }}

If you did not request this verification, please ignore this email.
`))

type SendEmailPasswordResetRequest struct {
	ProjectID         string
	EmailAddress      string
	PasswordResetCode string
}

func (s *Store) SendEmailPasswordReset(ctx context.Context, req *SendEmailPasswordResetRequest) error {
	projectID, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	qProject, err := s.q().GetProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	if qProject.CustomEmailPasswordReset {
		if err := s.SendWebhook(ctx, &SendWebhookRequest{
			ProjectID: req.ProjectID,
			EventType: "custom_email.password_reset",
			Payload: map[string]any{
				"type":              "custom_email.password_reset",
				"emailAddress":      req.EmailAddress,
				"passwordResetCode": req.PasswordResetCode,
			},
		}); err != nil {
			return fmt.Errorf("send webhook: %w", err)
		}
		return nil
	}

	var body bytes.Buffer
	if err := passwordResetEmailBodyTmpl.Execute(&body, struct {
		ProjectDisplayName string
		PasswordResetCode  string
	}{
		ProjectDisplayName: qProject.DisplayName,
		PasswordResetCode:  req.PasswordResetCode,
	}); err != nil {
		return fmt.Errorf("execute password reset email body template: %w", err)
	}

	if _, err := s.SES.SendEmail(ctx, &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(fmt.Sprintf("%s - Reset password", qProject.DisplayName)),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(body.String()),
					},
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{req.EmailAddress},
		},
		FromEmailAddress: aws.String(fmt.Sprintf("noreply@%s", qProject.EmailSendFromDomain)),
	}); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

var passwordResetEmailBodyTmpl = template.Must(template.New("passwordResetEmailBody").Parse(`Hello,

Someone has requested a password reset for your {{ .ProjectDisplayName }} account. If you did not request this, please ignore this email.

To continue logging in to {{ .ProjectDisplayName }}, please go back to the "Forgot password" page and enter this verification code:

{{ .PasswordResetCode }}

If you did not request this verification, please ignore this email.
`))

var userInviteEmailBodyTmpl = template.Must(template.New("userInviteEmailBodyTmpl").Parse(`Hello,

You have been invited to join {{ .OrganizationDisplayName }} in {{ .ProjectDisplayName }}.

You can accept this invite by signing up for {{ .ProjectDisplayName }}:

{{ .SignupLink }}
`))

type SendEmailUserInviteRequest struct {
	ProjectID    string
	UserInviteID string
}

func (s *Store) SendEmailUserInvite(ctx context.Context, req *SendEmailUserInviteRequest) error {
	projectID, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	qProject, err := s.q().GetProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	userInviteID, err := idformat.UserInvite.Parse(req.UserInviteID)
	if err != nil {
		return fmt.Errorf("parse user invite id: %w", err)
	}

	qUserInvite, err := s.q().GetUserInvite(ctx, userInviteID)
	if err != nil {
		return fmt.Errorf("get user invite: %w", err)
	}

	qOrganization, err := s.q().GetOrganization(ctx, qUserInvite.OrganizationID)
	if err != nil {
		return fmt.Errorf("get organization: %w", err)
	}

	if qProject.CustomEmailUserInvite {
		if err := s.SendWebhook(ctx, &SendWebhookRequest{
			ProjectID: req.ProjectID,
			EventType: "custom_email.user_invite",
			Payload: map[string]any{
				"type":         "custom_email.user_invite",
				"userInviteId": req.UserInviteID,
			},
		}); err != nil {
			return fmt.Errorf("send webhook: %w", err)
		}
		return nil
	}

	vaultDomain := qProject.VaultDomain
	if req.ProjectID == s.ConsoleProjectID {
		vaultDomain = s.ConsoleDomain
	}

	var body bytes.Buffer
	if err := userInviteEmailBodyTmpl.Execute(&body, struct {
		ProjectDisplayName      string
		OrganizationDisplayName string
		SignupLink              string
	}{
		ProjectDisplayName:      qProject.DisplayName,
		OrganizationDisplayName: qOrganization.DisplayName,
		SignupLink:              fmt.Sprintf("https://%s/signup", vaultDomain),
	}); err != nil {
		return fmt.Errorf("execute user invite email body template: %w", err)
	}

	if _, err := s.SES.SendEmail(ctx, &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(fmt.Sprintf("%s - You've been invited to join %s", qProject.DisplayName, qOrganization.DisplayName)),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(body.String()),
					},
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{qUserInvite.Email},
		},
		FromEmailAddress: aws.String(fmt.Sprintf("noreply@%s", qProject.EmailSendFromDomain)),
	}); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
