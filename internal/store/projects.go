package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/bcryptcost"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/store/queries"
	"golang.org/x/crypto/bcrypt"
)

type CreateProjectRequest struct {
	ConsoleProjectID string
	ProjectName      string
	OwnerUserEmail   string
	VaultDomain      string
	LogoURL          string
	DarkModeLogoURL  string
}

type CreateProjectResponse struct {
	ProjectID     string
	OwnerPassword string
}

func (s *Store) CreateProject(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	consoleProjectID, err := idformat.Project.Parse(req.ConsoleProjectID)
	if err != nil {
		return nil, fmt.Errorf("parse console project id: %w", err)
	}

	projectID := uuid.New()
	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                 uuid.New(),
		ProjectID:          consoleProjectID,
		DisplayName:        fmt.Sprintf("%s Backing Organization", idformat.Project.Format(projectID)),
		LogInWithGoogle:    false,
		LogInWithMicrosoft: false,
		LogInWithEmail:     true,
		LogInWithPassword:  true,
		ScimEnabled:        false,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	qProject, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                  projectID,
		OrganizationID:      &qOrganization.ID,
		DisplayName:         req.ProjectName,
		RedirectUri:         fmt.Sprintf("https://%s", req.VaultDomain),
		LogInWithGoogle:     false,
		LogInWithMicrosoft:  false,
		LogInWithEmail:      true,
		LogInWithPassword:   true,
		VaultDomain:         req.VaultDomain,
		EmailSendFromDomain: fmt.Sprintf("mail.%s", req.VaultDomain),
		CookieDomain:        req.VaultDomain,
	})
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	if _, err := q.CreateProjectTrustedDomain(ctx, queries.CreateProjectTrustedDomainParams{
		ID:        uuid.New(),
		ProjectID: qProject.ID,
		Domain:    req.VaultDomain,
	}); err != nil {
		return nil, fmt.Errorf("create project trusted domain: %w", err)
	}

	if _, err := q.CreateProjectUISettings(ctx, queries.CreateProjectUISettingsParams{
		ID:        uuid.New(),
		ProjectID: qProject.ID,
	}); err != nil {
		return nil, fmt.Errorf("create project ui settings: %w", err)
	}

	if req.LogoURL != "" {
		if _, err := q.UpdateProjectLogoURL(ctx, queries.UpdateProjectLogoURLParams{
			ProjectID: qProject.ID,
			LogoUrl:   &req.LogoURL,
		}); err != nil {
			return nil, fmt.Errorf("update project logo url: %w", err)
		}
	}

	if req.DarkModeLogoURL != "" {
		if _, err := q.UpdateProjectDarkModeLogoURL(ctx, queries.UpdateProjectDarkModeLogoURLParams{
			ProjectID:       qProject.ID,
			DarkModeLogoUrl: &req.DarkModeLogoURL,
		}); err != nil {
			return nil, fmt.Errorf("update project dark mode logo url: %w", err)
		}
	}

	// generate a random password for the bootstrap user
	var randomBytes [16]byte
	if _, err := rand.Read(randomBytes[:]); err != nil {
		panic(fmt.Errorf("read random bytes: %w", err))
	}

	ownerPassword := fmt.Sprintf("this_is_a_very_sensitive_password_%s", hex.EncodeToString(randomBytes[:]))
	ownerUserPasswordBcryptBytes, err := bcrypt.GenerateFromPassword([]byte(ownerPassword), bcryptcost.Cost)
	if err != nil {
		panic(fmt.Errorf("bcrypt bootstrap user password: %w", err))
	}

	// create the bootstrap user inside the console organization
	bootstrapUserPasswordBcrypt := string(ownerUserPasswordBcryptBytes)
	if _, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:             uuid.New(),
		OrganizationID: qOrganization.ID,
		Email:          req.OwnerUserEmail,
		IsOwner:        true,
		PasswordBcrypt: &bootstrapUserPasswordBcrypt,
	}); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// create session signing keys for the new project
	// Allow this key to be used for one year since the key rotation isn't implemented yet
	expiresAt := time.Now().Add(time.Hour * 24 * 365)

	// Generate a new symmetric key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	sskEncryptedBytes, err := s.sessionSigningKeyKMS.Encrypt(ctx, privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("encrypt session signing key: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, err
	}

	// Store the encrypted key in the database
	if _, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            qProject.ID,
		ExpireTime:           &expiresAt,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: sskEncryptedBytes,
	}); err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &CreateProjectResponse{
		ProjectID:     idformat.Project.Format(qProject.ID),
		OwnerPassword: ownerPassword,
	}, nil
}

type UpdateProjectRequest struct {
	ProjectID       string
	VaultDomain     string
	LogoURL         string
	DarkModeLogoURL string
}

func (s *Store) UpdateProject(ctx context.Context, req *UpdateProjectRequest) error {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}

	if req.VaultDomain != "" {
		if _, err := q.UpdateProject(ctx, queries.UpdateProjectParams{
			ID:          projectID,
			VaultDomain: req.VaultDomain,
		}); err != nil {
			return fmt.Errorf("update project: %w", err)
		}
	}

	if req.LogoURL != "" {
		if _, err := q.UpdateProjectLogoURL(ctx, queries.UpdateProjectLogoURLParams{
			ProjectID: projectID,
			LogoUrl:   &req.LogoURL,
		}); err != nil {
			return fmt.Errorf("update project logo url: %w", err)
		}
	}

	if req.DarkModeLogoURL != "" {
		if _, err := q.UpdateProjectDarkModeLogoURL(ctx, queries.UpdateProjectDarkModeLogoURLParams{
			ProjectID:       projectID,
			DarkModeLogoUrl: &req.DarkModeLogoURL,
		}); err != nil {
			return fmt.Errorf("update project dark mode logo url: %w", err)
		}
	}

	if err := commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
