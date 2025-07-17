package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stripe/stripe-go/v82"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"golang.org/x/net/publicsuffix"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var apiKeySecretTokenPrefixRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

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

	return &backendv1.GetProjectResponse{Project: s.parseProject(&project, qProjectTrustedDomains)}, nil
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
	if req.Project.GoogleOauthClientId != nil {
		updates.GoogleOauthClientID = req.Project.GoogleOauthClientId
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
	if req.Project.MicrosoftOauthClientId != nil {
		updates.MicrosoftOauthClientID = req.Project.MicrosoftOauthClientId
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

	updates.GithubOauthClientID = qProject.GithubOauthClientID
	if req.Project.GithubOauthClientId != nil {
		updates.GithubOauthClientID = req.Project.GithubOauthClientId
	}

	updates.GithubOauthClientSecretCiphertext = qProject.GithubOauthClientSecretCiphertext
	if req.Project.GithubOauthClientSecret != "" {
		encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               &s.githubOAuthClientSecretsKMSKeyID,
			Plaintext:           []byte(req.Project.GithubOauthClientSecret),
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		})
		if err != nil {
			return nil, fmt.Errorf("encrypt github oauth client secret: %w", err)
		}

		updates.GithubOauthClientSecretCiphertext = encryptRes.CiphertextBlob
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

	updates.LogInWithGithub = qProject.LogInWithGithub
	if req.Project.LogInWithGithub != nil {
		// todo: validate that github is configured?
		updates.LogInWithGithub = *req.Project.LogInWithGithub
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

	updates.LogInWithOidc = qProject.LogInWithOidc
	if req.Project.LogInWithOidc != nil {
		updates.LogInWithOidc = *req.Project.LogInWithOidc
	}

	updates.LogInWithAuthenticatorApp = qProject.LogInWithAuthenticatorApp
	if req.Project.LogInWithAuthenticatorApp != nil {
		updates.LogInWithAuthenticatorApp = *req.Project.LogInWithAuthenticatorApp
	}

	updates.LogInWithPasskey = qProject.LogInWithPasskey
	if req.Project.LogInWithPasskey != nil {
		updates.LogInWithPasskey = *req.Project.LogInWithPasskey
	}

	updates.ApiKeysEnabled = qProject.ApiKeysEnabled
	if req.Project.ApiKeysEnabled != nil {
		updates.ApiKeysEnabled = *req.Project.ApiKeysEnabled
	}

	updates.AuditLogsEnabled = qProject.AuditLogsEnabled
	if req.Project.AuditLogsEnabled != nil {
		updates.AuditLogsEnabled = *req.Project.AuditLogsEnabled
	}

	updates.ApiKeySecretTokenPrefix = qProject.ApiKeySecretTokenPrefix
	if req.Project.ApiKeySecretTokenPrefix != nil {
		if len(*req.Project.ApiKeySecretTokenPrefix) > 64 {
			return nil, apierror.NewFailedPreconditionError("api key secret token prefix must be no longer than 64 characters", fmt.Errorf("api key secret token prefix too long: %s", *req.Project.ApiKeySecretTokenPrefix))
		}

		if !apiKeySecretTokenPrefixRegex.MatchString(*req.Project.ApiKeySecretTokenPrefix) {
			return nil, apierror.NewFailedPreconditionError("api key secret token prefix must contain only lowercase letters, numbers, and underscores", fmt.Errorf("api key secret token prefix contains invalid characters: %s", *req.Project.ApiKeySecretTokenPrefix))
		}

		updates.ApiKeySecretTokenPrefix = req.Project.ApiKeySecretTokenPrefix
	}

	updates.CookieDomain = qProject.CookieDomain
	if req.Project.CookieDomain != "" {
		// only allow updates to cookie domain if the vault domain is custom
		defaultVaultDomain := fmt.Sprintf("%s.%s", strings.ReplaceAll(idformat.Project.Format(qProject.ID), "_", "-"), s.authAppsRootDomain)
		if qProject.VaultDomain == defaultVaultDomain {
			return nil, apierror.NewFailedPreconditionError("cannot update cookie domain unless vault domain is custom", nil)
		}

		// do not allow leading "." in cookie domain; we will automatically add
		// it in Set-Cookie headers
		if strings.HasPrefix(req.Project.CookieDomain, ".") {
			return nil, apierror.NewFailedPreconditionError("cookie domain must not start with '.'", nil)
		}

		// do not allow cookie domain to be from the public suffix list
		publicSuffix, _ := publicsuffix.PublicSuffix(req.Project.CookieDomain)
		if publicSuffix == req.Project.CookieDomain {
			return nil, apierror.NewFailedPreconditionError("cookie domain must not be public suffix", nil)
		}

		updates.CookieDomain = req.Project.CookieDomain
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

	if !qUpdatedProject.LogInWithGithub {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_github")
		if err := q.DisableProjectOrganizationsLogInWithGithub(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with Github: %w", err)
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

	if !qUpdatedProject.LogInWithOidc {
		slog.InfoContext(ctx, "disable_project_organizations_log_in_with_oidc")
		if err := q.DisableProjectOrganizationsLogInWithOIDC(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with OIDC: %w", err)
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

		for trustedDomain := range trustedDomains {
			domain := strings.Split(trustedDomain, ":")[0] // Remove port if present

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

	return &backendv1.UpdateProjectResponse{Project: s.parseProject(&qUpdatedProject, qProjectTrustedDomains)}, nil
}

func (s *Store) ConsoleCreateProject(ctx context.Context, req *backendv1.ConsoleCreateProjectRequest) (*backendv1.ConsoleCreateProjectResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if req.DisplayName == "" {
		return nil, apierror.NewFailedPreconditionError("project display name must be provided", nil)
	}

	if req.AppUrl == "" {
		return nil, apierror.NewFailedPreconditionError("project app URL must be provided", nil)
	}

	appUrl, err := url.Parse(req.AppUrl)
	if err != nil {
		return nil, fmt.Errorf("parse app url: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qConsoleProject, err := q.GetProjectByID(ctx, *s.dogfoodProjectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("console project not found", fmt.Errorf("get console project by id: %w", err))
		}
		return nil, fmt.Errorf("get console project by id: %w", err)
	}

	contextData := authn.GetContextData(ctx)
	userID, err := idformat.User.Parse(contextData.DogfoodSession.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	qUser, err := q.GetUser(ctx, queries.GetUserParams{
		ID:        userID,
		ProjectID: *s.dogfoodProjectID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user by id: %w", err))
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	projectID := uuid.New()
	organizationID := uuid.New()
	formattedNewProjectID := idformat.Project.Format(projectID)
	newProjectVaultDomain := fmt.Sprintf("%s.%s", strings.ReplaceAll(formattedNewProjectID, "_", "-"), s.authAppsRootDomain)

	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:          organizationID,
		ProjectID:   *s.dogfoodProjectID,
		DisplayName: fmt.Sprintf("%s Backing Organization", formattedNewProjectID),

		// same logic as in ordinary s.CreateOrganization, but mirroring the current Project's login methods
		LogInWithEmail:     qConsoleProject.LogInWithEmail,
		LogInWithGoogle:    qConsoleProject.LogInWithGoogle,
		LogInWithMicrosoft: qConsoleProject.LogInWithMicrosoft,
		LogInWithGithub:    qConsoleProject.LogInWithGithub,
		LogInWithPassword:  false,
		ScimEnabled:        false,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	_, err = q.CreateUser(ctx, queries.CreateUserParams{
		ID:              uuid.New(),
		OrganizationID:  qOrganization.ID,
		Email:           qUser.Email,
		GoogleUserID:    qUser.GoogleUserID,
		MicrosoftUserID: qUser.MicrosoftUserID,
		GithubUserID:    qUser.GithubUserID,
		IsOwner:         qUser.IsOwner,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	qNewProject, err := q.ConsoleCreateProject(ctx, queries.ConsoleCreateProjectParams{
		ID:                         projectID,
		OrganizationID:             refOrNil(qOrganization.ID),
		StripeCustomerID:           qProject.StripeCustomerID,
		DisplayName:                req.DisplayName,
		EntitledBackendApiKeys:     qProject.EntitledBackendApiKeys,
		EntitledCustomVaultDomains: qProject.EntitledCustomVaultDomains,
		LogInWithAuthenticatorApp:  qProject.LogInWithAuthenticatorApp,
		LogInWithEmail:             qProject.LogInWithEmail,
		LogInWithGithub:            qProject.LogInWithGithub,
		LogInWithGoogle:            qProject.LogInWithGoogle,
		LogInWithMicrosoft:         qProject.LogInWithMicrosoft,
		LogInWithOidc:              qProject.LogInWithOidc,
		LogInWithPasskey:           qProject.LogInWithPasskey,
		LogInWithPassword:          qProject.LogInWithPassword,
		LogInWithSaml:              qProject.LogInWithSaml,
		RedirectUri:                req.AppUrl,
		VaultDomain:                newProjectVaultDomain,
		CookieDomain:               newProjectVaultDomain,
		EmailSendFromDomain:        fmt.Sprintf("mail.%s", s.authAppsRootDomain),
	})
	if err != nil {
		return nil, fmt.Errorf("console create project: %w", err)
	}

	if _, err := q.CreateProjectUISettings(ctx, qNewProject.ID); err != nil {
		return nil, fmt.Errorf("create project ui settings: %w", err)
	}

	publicKey, privateKeyCiphertext, err := s.generateSessionSigningKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate prod session signing key: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	createTime := time.Now()
	expireTime := createTime.Add(time.Hour * 24 * 7)

	if _, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            qNewProject.ID,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: privateKeyCiphertext,
		CreateTime:           &createTime,
		ExpireTime:           &expireTime,
	}); err != nil {
		return nil, fmt.Errorf("create session signing key: %w", err)
	}

	// add the default set of trusted domains: the tesseral.app domain, and the
	// app url hostname
	if _, err := q.CreateProjectTrustedDomain(ctx, queries.CreateProjectTrustedDomainParams{
		ID:        uuid.New(),
		ProjectID: qNewProject.ID,
		Domain:    newProjectVaultDomain,
	}); err != nil {
		return nil, fmt.Errorf("create project trusted domain: %w", err)
	}

	if _, err := q.CreateProjectTrustedDomain(ctx, queries.CreateProjectTrustedDomainParams{
		ID:        uuid.New(),
		ProjectID: qNewProject.ID,
		Domain:    appUrl.Host,
	}); err != nil {
		return nil, fmt.Errorf("create project trusted domain: %w", err)
	}

	// Create a Svix application for the project to send webhooks to
	if _, err := s.createProjectWebhookSettings(ctx, q, qProject); err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}

	// Add the new project to the Stripe customer metadata
	if err := s.upsertStripeCustomerProjectID(ctx, idformat.Project.Format(qNewProject.ID), qProject.StripeCustomerID); err != nil {
		return nil, fmt.Errorf("upsert stripe customer project id: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.ConsoleCreateProjectResponse{
		Project:        s.parseProject(&qNewProject, []queries.ProjectTrustedDomain{}),
		OrganizationId: idformat.Organization.Format(qOrganization.ID),
	}, nil
}

func (s *Store) generateSessionSigningKey(ctx context.Context) (*ecdsa.PublicKey, []byte, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the symmetric key with the KMS
	sskEncryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
		Plaintext:           privateKeyBytes,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("encrypt session signing key: %w", err)
	}

	return privateKey.Public().(*ecdsa.PublicKey), sskEncryptOutput.CiphertextBlob, nil
}

func (s *Store) upsertStripeCustomerProjectID(ctx context.Context, projectID string, stripeCustomerID *string) error {
	if stripeCustomerID == nil || *stripeCustomerID == "" {
		return fmt.Errorf("stripe customer ID is required to upsert project ID")
	}

	customer, err := s.stripe.Customers.Get(*stripeCustomerID, &stripe.CustomerParams{})
	if err != nil {
		return fmt.Errorf("get stripe customer: %w", err)
	}

	var projectIDs []string
	var projectIDExists bool
	if customer.Metadata != nil {
		if projectIDStr, ok := customer.Metadata["project_id"]; ok {
			storedProjectIDs := strings.Split(projectIDStr, ",")

			// Looping over all stored project IDs to prevent duplicates
			for _, id := range storedProjectIDs {
				if id == projectID {
					projectIDExists = true
				}
				projectIDs = append(projectIDs, id)
			}
		}
	}
	// If the project ID doesn't exist, add it to the list
	if !projectIDExists {
		projectIDs = append(projectIDs, projectID)
	}

	if _, err := s.stripe.Customers.Update(*stripeCustomerID, &stripe.CustomerParams{
		Metadata: map[string]string{
			"project_id": strings.Join(projectIDs, ","),
		},
	}); err != nil {
		return fmt.Errorf("update stripe customer metadata: %w", err)
	}

	return nil

}

func (s *Store) parseProject(qProject *queries.Project, qProjectTrustedDomains []queries.ProjectTrustedDomain) *backendv1.Project {
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
		LogInWithGithub:            &qProject.LogInWithGithub,
		LogInWithEmail:             &qProject.LogInWithEmail,
		LogInWithPassword:          &qProject.LogInWithPassword,
		LogInWithSaml:              &qProject.LogInWithSaml,
		LogInWithOidc:              &qProject.LogInWithOidc,
		LogInWithAuthenticatorApp:  &qProject.LogInWithAuthenticatorApp,
		LogInWithPasskey:           &qProject.LogInWithPasskey,
		GoogleOauthClientId:        qProject.GoogleOauthClientID,
		GoogleOauthClientSecret:    "", // intentionally left blank
		MicrosoftOauthClientId:     qProject.MicrosoftOauthClientID,
		MicrosoftOauthClientSecret: "", // intentionally left blank
		GithubOauthClientId:        qProject.GithubOauthClientID,
		GithubOauthClientSecret:    "", // intentionally left blank
		VaultDomain:                qProject.VaultDomain,
		VaultDomainCustom:          qProject.VaultDomain != fmt.Sprintf("%s.%s", strings.ReplaceAll(idformat.Project.Format(qProject.ID), "_", "-"), s.authAppsRootDomain),
		TrustedDomains:             trustedDomains,
		CookieDomain:               qProject.CookieDomain,
		EmailSendFromDomain:        qProject.EmailSendFromDomain,
		ApiKeysEnabled:             &qProject.ApiKeysEnabled,
		ApiKeySecretTokenPrefix:    qProject.ApiKeySecretTokenPrefix,
		AuditLogsEnabled:           refOrNil(qProject.AuditLogsEnabled),
	}
}
