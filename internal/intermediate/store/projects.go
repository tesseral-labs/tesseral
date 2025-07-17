package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateProject(ctx context.Context, req *intermediatev1.CreateProjectRequest) (*intermediatev1.CreateProjectResponse, error) {
	if authn.ProjectID(ctx) != *s.dogfoodProjectID {
		return nil, apierror.NewPermissionDeniedError("cannot create a project", fmt.Errorf("create project attempted on non-dogfood project: %s", authn.ProjectID(ctx)))
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	if !intermediateSession.EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", nil)
	}

	// create this ahead of time so we can use it in the display name and auth domain
	newProjectID := uuid.New()
	formattedNewProjectID := idformat.Project.Format(newProjectID)
	newProjectVaultDomain := fmt.Sprintf("%s.%s", strings.ReplaceAll(formattedNewProjectID, "_", "-"), s.authAppsRootDomain)

	// create a new organization under the dogfood project, accepting the same
	// primary login method used to get to this point
	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                 uuid.New(),
		DisplayName:        fmt.Sprintf("%s Backing Organization", formattedNewProjectID),
		ProjectID:          *s.dogfoodProjectID,
		LogInWithEmail:     intermediateSession.PrimaryAuthFactor == intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_EMAIL,
		LogInWithGoogle:    intermediateSession.PrimaryAuthFactor == intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_GOOGLE,
		LogInWithMicrosoft: intermediateSession.PrimaryAuthFactor == intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_MICROSOFT,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	// reflect the google hosted domain from the intermediate session if it exists
	if intermediateSession.GoogleHostedDomain != "" {
		if _, err := q.CreateOrganizationGoogleHostedDomain(ctx, queries.CreateOrganizationGoogleHostedDomainParams{
			ID:                 uuid.New(),
			OrganizationID:     qOrganization.ID,
			GoogleHostedDomain: intermediateSession.GoogleHostedDomain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}

	// reflect the microsoft tenant id from the intermediate session if it exists
	if intermediateSession.MicrosoftTenantId != "" {
		if _, err := q.CreateOrganizationMicrosoftTenantID(ctx, queries.CreateOrganizationMicrosoftTenantIDParams{
			ID:                uuid.New(),
			OrganizationID:    qOrganization.ID,
			MicrosoftTenantID: intermediateSession.MicrosoftTenantId,
		}); err != nil {
			return nil, fmt.Errorf("create organization microsoft tenant id: %w", err)
		}
	}

	// create a new user invite for the intermediate session user
	if _, err := q.CreateUserInvite(ctx, queries.CreateUserInviteParams{
		ID:             uuid.New(),
		OrganizationID: qOrganization.ID,
		Email:          intermediateSession.Email,
		IsOwner:        true,
	}); err != nil {
		return nil, fmt.Errorf("create user invite: %w", err)
	}

	// create a new project backed by the new organization; the login methods
	// here are only those that can work out of the box, without further
	// configuration by the user
	qProject, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                  newProjectID,
		RedirectUri:         req.RedirectUri,
		OrganizationID:      &qOrganization.ID,
		VaultDomain:         newProjectVaultDomain,
		EmailSendFromDomain: fmt.Sprintf("mail.%s", s.authAppsRootDomain),
		DisplayName:         req.DisplayName,
		LogInWithEmail:      true,
		LogInWithGoogle:     false,
		LogInWithMicrosoft:  false,
		LogInWithPassword:   false,
		LogInWithSaml:       false,
	})
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	if _, err := q.CreateProjectUISettings(ctx, qProject.ID); err != nil {
		return nil, fmt.Errorf("create project ui settings: %w", err)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	sskEncryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
		Plaintext:           privateKeyBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("encrypt session signing key: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	createTime := time.Now()
	expireTime := createTime.Add(time.Hour * 24 * 7)

	if _, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            qProject.ID,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: sskEncryptOutput.CiphertextBlob,
		CreateTime:           &createTime,
		ExpireTime:           &expireTime,
	}); err != nil {
		return nil, fmt.Errorf("create session signing key: %w", err)
	}

	if err = commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.CreateProjectResponse{
		Project: parseProject(qProject),
	}, nil
}

func (s *Store) OnboardingCreateProjects(ctx context.Context, req *intermediatev1.OnboardingCreateProjectsRequest) (*intermediatev1.OnboardingCreateProjectsResponse, error) {
	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qIntermediateSession, err := s.q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	sandboxPublicKey, sandboxPrivateKeyCiphertext, err := s.onboardingGenerateSessionSigningKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate prod session signing key: %w", err)
	}

	qSandboxProjectID := uuid.New()

	slog.InfoContext(ctx, "stripe_client_exists", "stripe_client_exists", s.stripeClient != nil)

	var stripeCustomerID *string
	if s.stripeClient != nil {
		stripeCustomer, err := s.stripeClient.Customers.New(&stripe.CustomerParams{
			Name:  &req.DisplayName,
			Email: qIntermediateSession.Email,
			Metadata: map[string]string{
				"tesseral_project_ids": idformat.Project.Format(qSandboxProjectID),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("create stripe customer: %w", err)
		}

		slog.InfoContext(ctx, "create_stripe_customer", "stripe_customer_id", stripeCustomer.ID)
		stripeCustomerID = &stripeCustomer.ID
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qSandboxUser, err := s.createProjectForCurrentUser(ctx, q, &qIntermediateSession, createProjectForCurrentUserArgs{
		ProjectID:                          qSandboxProjectID,
		StripeCustomerID:                   stripeCustomerID,
		RedirectURI:                        req.AppUrl,
		DisplayName:                        req.DisplayName,
		SessionSigningPublicKey:            sandboxPublicKey,
		SessionSigningPrivateKeyCiphertext: sandboxPrivateKeyCiphertext,
	})
	if err != nil {
		return nil, fmt.Errorf("create sandbox project for current user: %w", err)
	}

	slog.InfoContext(ctx, "created_project",
		"display_name", req.DisplayName,
		"sandbox_project_id", idformat.Project.Format(qSandboxProjectID))

	expireTime := time.Now().Add(sessionDuration)

	// Create a new session for the user
	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])
	if _, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expireTime,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qSandboxUser.ID,
		PrimaryAuthFactor:  *qIntermediateSession.PrimaryAuthFactor,
	}); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	// revoke the intermediate session
	if _, err := q.RevokeIntermediateSession(ctx, qIntermediateSession.ID); err != nil {
		return nil, fmt.Errorf("revoke intermediate session: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.OnboardingCreateProjectsResponse{
		AccessToken:  "", // populated in service
		RefreshToken: idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func (s *Store) onboardingGenerateSessionSigningKey(ctx context.Context) (*ecdsa.PublicKey, []byte, error) {
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

type createProjectForCurrentUserArgs struct {
	ProjectID                          uuid.UUID
	StripeCustomerID                   *string
	RedirectURI                        string
	DisplayName                        string
	SessionSigningPublicKey            *ecdsa.PublicKey
	SessionSigningPrivateKeyCiphertext []byte
}

func (s *Store) createProjectForCurrentUser(ctx context.Context, q *queries.Queries, qIntermediateSession *queries.IntermediateSession, args createProjectForCurrentUserArgs) (*queries.User, error) {
	formattedNewProjectID := idformat.Project.Format(args.ProjectID)
	newProjectVaultDomain := fmt.Sprintf("%s.%s", strings.ReplaceAll(formattedNewProjectID, "_", "-"), s.authAppsRootDomain)

	redirectURI, err := url.Parse(args.RedirectURI)
	if err != nil {
		return nil, fmt.Errorf("parse redirect uri: %w", err)
	}

	qDogfoodProject, err := q.GetProjectByID(ctx, *s.dogfoodProjectID)
	if err != nil {
		return nil, fmt.Errorf("get dogfood project by id: %w", err)
	}

	// create a new organization under the dogfood project, accepting the same
	// primary login method used to get to this point
	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:          uuid.New(),
		DisplayName: fmt.Sprintf("%s Backing Organization", formattedNewProjectID),
		ProjectID:   *s.dogfoodProjectID,

		// same logic as in ordinary s.CreateOrganization, but against dogfood project
		LogInWithEmail:     qDogfoodProject.LogInWithEmail,
		LogInWithGoogle:    qDogfoodProject.LogInWithGoogle,
		LogInWithMicrosoft: qDogfoodProject.LogInWithMicrosoft,
		LogInWithGithub:    qDogfoodProject.LogInWithGithub,
		LogInWithPassword:  qDogfoodProject.LogInWithPassword,
		ScimEnabled:        false,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	// reflect the google hosted domain from the intermediate session if it exists
	if qIntermediateSession.GoogleHostedDomain != nil {
		if _, err := q.CreateOrganizationGoogleHostedDomain(ctx, queries.CreateOrganizationGoogleHostedDomainParams{
			ID:                 uuid.New(),
			OrganizationID:     qOrganization.ID,
			GoogleHostedDomain: *qIntermediateSession.GoogleHostedDomain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}

	// reflect the microsoft tenant id from the intermediate session if it exists
	if qIntermediateSession.MicrosoftTenantID != nil {
		if _, err := q.CreateOrganizationMicrosoftTenantID(ctx, queries.CreateOrganizationMicrosoftTenantIDParams{
			ID:                uuid.New(),
			OrganizationID:    qOrganization.ID,
			MicrosoftTenantID: *qIntermediateSession.MicrosoftTenantID,
		}); err != nil {
			return nil, fmt.Errorf("create organization microsoft tenant id: %w", err)
		}
	}

	// create a user from the intermediate session
	qUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:                uuid.New(),
		OrganizationID:    qOrganization.ID,
		Email:             *qIntermediateSession.Email,
		DisplayName:       qIntermediateSession.UserDisplayName,
		ProfilePictureUrl: qIntermediateSession.ProfilePictureUrl,
		GoogleUserID:      qIntermediateSession.GoogleUserID,
		MicrosoftUserID:   qIntermediateSession.MicrosoftUserID,
		IsOwner:           true,
		PasswordBcrypt:    qIntermediateSession.NewUserPasswordBcrypt,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// create a new project backed by the new organization; the login methods
	// here are only those that can work out of the box, without further
	// configuration by the user
	qProject, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                  args.ProjectID,
		StripeCustomerID:    args.StripeCustomerID,
		RedirectUri:         args.RedirectURI,
		OrganizationID:      &qOrganization.ID,
		VaultDomain:         newProjectVaultDomain,
		CookieDomain:        newProjectVaultDomain,
		EmailSendFromDomain: fmt.Sprintf("mail.%s", s.authAppsRootDomain),
		DisplayName:         args.DisplayName,
		LogInWithEmail:      true,
		LogInWithGoogle:     false,
		LogInWithMicrosoft:  false,
		LogInWithPassword:   false,
		LogInWithSaml:       false,
	})
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	if _, err := q.CreateProjectUISettings(ctx, qProject.ID); err != nil {
		return nil, fmt.Errorf("create project ui settings: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(args.SessionSigningPublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	createTime := time.Now()
	expireTime := createTime.Add(time.Hour * 24 * 7)

	if _, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            qProject.ID,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: args.SessionSigningPrivateKeyCiphertext,
		CreateTime:           &createTime,
		ExpireTime:           &expireTime,
	}); err != nil {
		return nil, fmt.Errorf("create session signing key: %w", err)
	}

	// add the default set of trusted domains: the tesseral.app domain, and the
	// redirect URI hostname
	if _, err := q.CreateProjectTrustedDomain(ctx, queries.CreateProjectTrustedDomainParams{
		ID:        uuid.New(),
		ProjectID: qProject.ID,
		Domain:    newProjectVaultDomain,
	}); err != nil {
		return nil, fmt.Errorf("create project trusted domain: %w", err)
	}

	if _, err := q.CreateProjectTrustedDomain(ctx, queries.CreateProjectTrustedDomainParams{
		ID:        uuid.New(),
		ProjectID: qProject.ID,
		Domain:    redirectURI.Host,
	}); err != nil {
		return nil, fmt.Errorf("create project trusted domain: %w", err)
	}

	// Create a Svix application for the project to send webhooks to
	if _, err := s.createProjectWebhookSettings(ctx, q, qProject); err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}

	// Create a Publishable Key for the Project
	if _, err := q.CreatePublishableKey(ctx, queries.CreatePublishableKeyParams{
		ID:          uuid.New(),
		ProjectID:   args.ProjectID,
		DisplayName: "Default Publishable Key",
		DevMode:     true,
	}); err != nil {
		return nil, fmt.Errorf("create publishable key: %w", err)
	}

	return &qUser, nil
}

func parseProject(qProject queries.Project) *intermediatev1.Project {
	return &intermediatev1.Project{
		Id:             idformat.Project.Format(qProject.ID),
		OrganizationId: idformat.Organization.Format(*qProject.OrganizationID),
		CreateTime:     timestamppb.New(*qProject.CreateTime),
		UpdateTime:     timestamppb.New(*qProject.UpdateTime),
		DisplayName:    qProject.DisplayName,
		VaultDomain:    qProject.VaultDomain,
	}
}
