package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetProject(ctx context.Context, req *backendv1.GetProjectRequest) (*backendv1.GetProjectResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	project, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, err
	}

	qProjectTrustedDomains, err := q.GetProjectTrustedDomains(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project trusted domains: %w", err)
	}

	return &backendv1.GetProjectResponse{Project: parseProject(&project, qProjectTrustedDomains)}, nil
}

func (s *Store) DisableProjectLogins(ctx context.Context, req *backendv1.DisableProjectLoginsRequest) (*backendv1.DisableProjectLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.DisableProjectLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("lockout project: %w", err)
	}

	if err := q.RevokeAllProjectSessions(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("revoke all project sessions: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DisableProjectLoginsResponse{}, nil
}

func (s *Store) EnableProjectLogins(ctx context.Context, req *backendv1.EnableProjectLoginsRequest) (*backendv1.EnableProjectLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.EnableProjectLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("unlock project: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.EnableProjectLoginsResponse{}, nil
}

func (s *Store) UpdateProject(ctx context.Context, req *backendv1.UpdateProjectRequest) (*backendv1.UpdateProjectResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	// fetch project outside a transaction, so that we can carry out KMS
	// operations; we can live with possibility of conflicting concurrent writes
	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	updates := queries.UpdateProjectParams{
		ID: qProject.ID,
	}

	updates.DisplayName = qProject.DisplayName
	if req.Project.DisplayName != "" {
		updates.DisplayName = req.Project.DisplayName
	}

	updates.GoogleOauthClientID = qProject.GoogleOauthClientID
	if req.Project.GoogleOauthClientId != "" {
		updates.GoogleOauthClientID = &req.Project.GoogleOauthClientId
	}

	updates.GoogleOauthClientSecretCiphertext = qProject.GoogleOauthClientSecretCiphertext
	if req.Project.GoogleOauthClientSecret != "" {
		encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               &s.googleOAuthClientSecretsKMSKeyID,
			Plaintext:           []byte(req.Project.GoogleOauthClientSecret),
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		})
		if err != nil {
			return nil, fmt.Errorf("encrypt google oauth client secret: %w", err)
		}

		updates.GoogleOauthClientSecretCiphertext = encryptRes.CiphertextBlob
	}

	updates.MicrosoftOauthClientID = qProject.MicrosoftOauthClientID
	if req.Project.MicrosoftOauthClientId != "" {
		updates.MicrosoftOauthClientID = &req.Project.MicrosoftOauthClientId
	}

	updates.MicrosoftOauthClientSecretCiphertext = qProject.MicrosoftOauthClientSecretCiphertext
	if req.Project.MicrosoftOauthClientSecret != "" {
		encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               &s.microsoftOAuthClientSecretsKMSKeyID,
			Plaintext:           []byte(req.Project.MicrosoftOauthClientSecret),
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		})
		if err != nil {
			return nil, fmt.Errorf("encrypt microsoft oauth client secret: %w", err)
		}

		updates.MicrosoftOauthClientSecretCiphertext = encryptRes.CiphertextBlob
	}

	updates.LogInWithGoogle = qProject.LogInWithGoogle
	if req.Project.LogInWithGoogle != nil {
		// todo: validate that google is configured?
		updates.LogInWithGoogle = *req.Project.LogInWithGoogle
	}

	updates.LogInWithMicrosoft = qProject.LogInWithMicrosoft
	if req.Project.LogInWithMicrosoft != nil {
		// todo: validate that microsoft is configured?
		updates.LogInWithMicrosoft = *req.Project.LogInWithMicrosoft
	}

	updates.LogInWithEmail = qProject.LogInWithEmail
	if req.Project.LogInWithEmail != nil {
		updates.LogInWithEmail = *req.Project.LogInWithEmail
	}

	updates.LogInWithPassword = qProject.LogInWithPassword
	if req.Project.LogInWithPassword != nil {
		updates.LogInWithPassword = *req.Project.LogInWithPassword
	}

	updates.LogInWithSaml = qProject.LogInWithSaml
	if req.Project.LogInWithSaml != nil {
		updates.LogInWithSaml = *req.Project.LogInWithSaml
	}

	updates.LogInWithAuthenticatorApp = qProject.LogInWithAuthenticatorApp
	if req.Project.LogInWithAuthenticatorApp != nil {
		updates.LogInWithAuthenticatorApp = *req.Project.LogInWithAuthenticatorApp
	}

	updates.LogInWithPasskey = qProject.LogInWithPasskey
	if req.Project.LogInWithPasskey != nil {
		updates.LogInWithPasskey = *req.Project.LogInWithPasskey
	}

	updates.RedirectUri = qProject.RedirectUri
	if req.Project.RedirectUri != "" {
		updates.RedirectUri = req.Project.RedirectUri
	}

	updates.AfterLoginRedirectUri = qProject.AfterLoginRedirectUri
	if req.Project.AfterLoginRedirectUri != nil {
		updates.AfterLoginRedirectUri = req.Project.AfterLoginRedirectUri
	}

	updates.AfterSignupRedirectUri = qProject.AfterSignupRedirectUri
	if req.Project.AfterSignupRedirectUri != nil {
		updates.AfterSignupRedirectUri = req.Project.AfterSignupRedirectUri
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qUpdatedProject, err := q.UpdateProject(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	if !qUpdatedProject.LogInWithGoogle {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_google")
		if err := q.DisableProjectOrganizationsLogInWithGoogle(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with Google: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithMicrosoft {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_microsoft")
		if err := q.DisableProjectOrganizationsLogInWithMicrosoft(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with Microsoft: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithEmail {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_email")
		if err := q.DisableProjectOrganizationsLogInWithEmail(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with email: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithPassword {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_password")
		if err := q.DisableProjectOrganizationsLogInWithPassword(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with password: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithSaml {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_saml")
		if err := q.DisableProjectOrganizationsLogInWithSAML(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with SAML: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithAuthenticatorApp {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_authenticator_app")
		if err := q.DisableProjectOrganizationsLogInWithAuthenticatorApp(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with authenticator app: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithPasskey {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_passkey")
		if err := q.DisableProjectOrganizationsLogInWithPasskey(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with passkey: %w", err)
		}
	}

	// only update project trusted domains if mentioned in request
	if len(req.Project.TrustedDomains) > 0 {
		// always include the default vault domain (project-xxx.tesseral.app)
		// and the current vault domain (e.g. auth.company.com) in the set of
		// trusted domains
		trustedDomains := map[string]struct{}{
			qUpdatedProject.VaultDomain: {},
			fmt.Sprintf("%s.%s", strings.ReplaceAll(idformat.Project.Format(qUpdatedProject.ID), "_", "-"), s.authAppsRootDomain): {},
		}
		for _, domain := range req.Project.TrustedDomains {
			trustedDomains[domain] = struct{}{}
		}

		if err := q.DeleteProjectTrustedDomainsByProjectID(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("delete project trusted domains by project id: %w", err)
		}

		for domain := range trustedDomains {
			if _, err := q.CreateProjectTrustedDomain(ctx, queries.CreateProjectTrustedDomainParams{
				ID:        uuid.New(),
				ProjectID: authn.ProjectID(ctx),
				Domain:    domain,
			}); err != nil {
				return nil, fmt.Errorf("create project passkey rp id: %w", err)
			}
		}
	}

	qProjectTrustedDomains, err := q.GetProjectTrustedDomains(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project trusted domains: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateProjectResponse{Project: parseProject(&qUpdatedProject, qProjectTrustedDomains)}, nil
}

func parseProject(qProject *queries.Project, qProjectTrustedDomains []queries.ProjectTrustedDomain) *backendv1.Project {
	// sanity check
	for _, qProjectTrustedDomain := range qProjectTrustedDomains {
		if qProjectTrustedDomain.ProjectID != qProject.ID {
			panic(fmt.Errorf("project trusted domain project id mismatch: %s != %s", qProjectTrustedDomain.ProjectID, qProject.ID))
		}
	}

	var trustedDomains []string
	for _, qProjectTrustedDomain := range qProjectTrustedDomains {
		trustedDomains = append(trustedDomains, qProjectTrustedDomain.Domain)
	}

	return &backendv1.Project{
		Id:                         idformat.Project.Format(qProject.ID),
		DisplayName:                qProject.DisplayName,
		CreateTime:                 timestamppb.New(*qProject.CreateTime),
		UpdateTime:                 timestamppb.New(*qProject.UpdateTime),
		LogInWithGoogle:            &qProject.LogInWithGoogle,
		LogInWithMicrosoft:         &qProject.LogInWithMicrosoft,
		LogInWithEmail:             &qProject.LogInWithEmail,
		LogInWithPassword:          &qProject.LogInWithPassword,
		LogInWithSaml:              &qProject.LogInWithSaml,
		LogInWithAuthenticatorApp:  &qProject.LogInWithAuthenticatorApp,
		LogInWithPasskey:           &qProject.LogInWithPasskey,
		GoogleOauthClientId:        derefOrEmpty(qProject.GoogleOauthClientID),
		GoogleOauthClientSecret:    "", // intentionally left blank
		MicrosoftOauthClientId:     derefOrEmpty(qProject.MicrosoftOauthClientID),
		MicrosoftOauthClientSecret: "", // intentionally left blank
		VaultDomain:                qProject.VaultDomain,
		TrustedDomains:             trustedDomains,
		RedirectUri:                qProject.RedirectUri,
		AfterLoginRedirectUri:      qProject.AfterLoginRedirectUri,
		AfterSignupRedirectUri:     qProject.AfterSignupRedirectUri,
		EmailSendFromDomain:        qProject.EmailSendFromDomain,
	}
}
