package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// defaultEmailQuotaDaily is the default number of emails a project may send per
// day.
var defaultEmailQuotaDaily int32 = 1000

func (s *Store) ListUserInvites(ctx context.Context, req *frontendv1.ListUserInvitesRequest) (*frontendv1.ListUserInvitesResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qUserInvites, err := q.ListUserInvites(ctx, queries.ListUserInvitesParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list user invites: %w", err)
	}

	var userInvites []*frontendv1.UserInvite
	for _, qUserInvite := range qUserInvites {
		userInvites = append(userInvites, parseUserInvite(qUserInvite))
	}

	var nextPageToken string
	if len(userInvites) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qUserInvites[limit].ID)
		userInvites = userInvites[:limit]
	}

	return &frontendv1.ListUserInvitesResponse{
		UserInvites:   userInvites,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetUserInvite(ctx context.Context, req *frontendv1.GetUserInviteRequest) (*frontendv1.GetUserInviteResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userInviteID, err := idformat.UserInvite.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse user invite id: %w", err)
	}

	qInvite, err := q.GetUserInvite(ctx, queries.GetUserInviteParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userInviteID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user invite not found", fmt.Errorf("get user invite: %w", err))
		}

		return nil, fmt.Errorf("get user invite: %w", err)
	}

	return &frontendv1.GetUserInviteResponse{UserInvite: parseUserInvite(qInvite)}, nil
}

func (s *Store) CreateUserInvite(ctx context.Context, req *frontendv1.CreateUserInviteRequest) (*frontendv1.CreateUserInviteResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	// As a sanity measure, prevent invites to existing users with the same
	// email. Otherwise, those invites will remain present even after the user
	// is deleted, and so someone with that email address could redeem the
	// invite afterword.
	emailTaken, err := q.ExistsUserWithEmail(ctx, queries.ExistsUserWithEmailParams{
		OrganizationID: authn.OrganizationID(ctx),
		Email:          req.UserInvite.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("exists user with email: %w", err)
	}

	if emailTaken {
		return nil, apierror.NewFailedPreconditionError("a user with that email already exists", nil)
	}

	qUserInvite, err := q.CreateUserInvite(ctx, queries.CreateUserInviteParams{
		ID:             uuid.New(),
		OrganizationID: authn.OrganizationID(ctx),
		Email:          req.UserInvite.Email,
		IsOwner:        req.UserInvite.Owner,
	})
	if err != nil {
		return nil, fmt.Errorf("create user invite: %w", err)
	}

	qEmailDailyQuotaUsage, err := q.IncrementProjectEmailDailyQuotaUsage(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("increment project email daily quota usage: %w", err)
	}

	emailQuotaDaily := defaultEmailQuotaDaily
	if qProject.EmailQuotaDaily != nil {
		emailQuotaDaily = *qProject.EmailQuotaDaily
	}

	slog.InfoContext(ctx, "email_daily_quota_usage", "usage", qEmailDailyQuotaUsage.QuotaUsage, "quota", qProject.EmailQuotaDaily)

	if qEmailDailyQuotaUsage.QuotaUsage > emailQuotaDaily {
		slog.InfoContext(ctx, "email_daily_quota_exceeded")
		return nil, apierror.NewFailedPreconditionError("email daily quota exceeded", fmt.Errorf("email daily quota exceeded"))
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	if req.SendEmail {
		if err := s.sendUserInviteEmail(ctx, req.UserInvite.Email, qOrg.DisplayName); err != nil {
			return nil, fmt.Errorf("send user invite email: %w", err)
		}
	}

	userInvite := parseUserInvite(qUserInvite)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.user_invites.create",
		EventDetails: map[string]any{
			"userInvite": userInvite,
		},
		ResourceType: queries.AuditLogEventResourceTypeUserInvite,
		ResourceID:   &qUserInvite.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	return &frontendv1.CreateUserInviteResponse{UserInvite: userInvite}, nil
}

func (s *Store) DeleteUserInvite(ctx context.Context, req *frontendv1.DeleteUserInviteRequest) (*frontendv1.DeleteUserInviteResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userInviteID, err := idformat.UserInvite.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse user invite id: %w", err)
	}

	qUserInvite, err := q.GetUserInvite(ctx, queries.GetUserInviteParams{
		ID:             userInviteID,
		OrganizationID: authn.OrganizationID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user invite not found", fmt.Errorf("get user invite: %w", err))
		}

		return nil, fmt.Errorf("get user invite: %w", err)
	}

	if err := q.DeleteUserInvite(ctx, userInviteID); err != nil {
		return nil, fmt.Errorf("delete user invite: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	userInvite := parseUserInvite(qUserInvite)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.user_invites.delete",
		EventDetails: map[string]any{
			"userInvite": userInvite,
		},
		ResourceType: queries.AuditLogEventResourceTypeUserInvite,
		ResourceID:   &qUserInvite.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	return &frontendv1.DeleteUserInviteResponse{}, nil
}

func parseUserInvite(qUserInvite queries.UserInvite) *frontendv1.UserInvite {
	return &frontendv1.UserInvite{
		Id:         idformat.UserInvite.Format(qUserInvite.ID),
		CreateTime: timestamppb.New(*qUserInvite.CreateTime),
		UpdateTime: timestamppb.New(*qUserInvite.UpdateTime),
		Email:      qUserInvite.Email,
		Owner:      qUserInvite.IsOwner,
	}
}

var userInviteEmailBodyTmpl = template.Must(template.New("userInviteEmailBodyTmpl").Parse(`Hello,

You have been invited to join {{ .OrganizationDisplayName }} in {{ .ProjectDisplayName }}.

You can accept this invite by signing up for {{ .ProjectDisplayName }}:

{{ .SignupLink }}
`))

func (s *Store) sendUserInviteEmail(ctx context.Context, toAddress string, organizationDisplayName string) error {
	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	subject := fmt.Sprintf("%s - You've been invited to join %s", qProject.DisplayName, organizationDisplayName)

	vaultDomain := qProject.VaultDomain
	if authn.ProjectID(ctx) == *s.dogfoodProjectID {
		vaultDomain = s.consoleDomain
	}

	var body bytes.Buffer
	if err := userInviteEmailBodyTmpl.Execute(&body, struct {
		ProjectDisplayName      string
		OrganizationDisplayName string
		SignupLink              string
	}{
		ProjectDisplayName:      qProject.DisplayName,
		OrganizationDisplayName: organizationDisplayName,
		SignupLink:              fmt.Sprintf("https://%s/signup", vaultDomain),
	}); err != nil {
		return fmt.Errorf("execute email verification email body template: %w", err)
	}

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
