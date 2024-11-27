package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type CreateDogfoodProjectResponse struct {
	DogfoodProjectID                   string
	BootstrapUserEmail                 string
	BootstrapUserVerySensitivePassword string
	SessionSigningKeyID             		string
	IntermediateSessionSigningKeyID 		string
}

// CreateDogfoodProject creates the dogfood project.
func (s *Store) CreateDogfoodProject(ctx context.Context) (*CreateDogfoodProjectResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// refuse to proceed if the database has any existing projects
	projectCount, err := q.CountAllProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("count projects: %w", err)
	}

	if projectCount != 0 {
		return nil, fmt.Errorf("project count is not zero: %d", projectCount)
	}

	// directly create project and organization that cyclically refer to each
	// other
	dogfoodProjectID := uuid.New()
	dogfoodOrganizationID := uuid.New()
	formattedDogfoodProjectID := idformat.Project.Format(dogfoodProjectID)

	if _, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                       dogfoodProjectID,
		OrganizationID:           nil, // will populate after creating org
		LogInWithPasswordEnabled: true,
	}); err != nil {
		return nil, fmt.Errorf("create dogfood project: %w", err)
	}

	if _, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:          dogfoodOrganizationID,
		ProjectID:   dogfoodProjectID,
		DisplayName: "OpenAuth",
	}); err != nil {
		return nil, fmt.Errorf("create dogfood organization: %w", err)
	}

	// manually link project to organization
	if _, err := q.UpdateProjectOrganizationID(ctx, queries.UpdateProjectOrganizationIDParams{
		ID:             dogfoodProjectID,
		OrganizationID: &dogfoodOrganizationID,
	}); err != nil {
		return nil, fmt.Errorf("update dogfood project organization: %w", err)
	}

	// generate a random password for the bootstrap user
	var randomBytes [16]byte
	if _, err := rand.Read(randomBytes[:]); err != nil {
		panic(fmt.Errorf("read random bytes: %w", err))
	}

	bootstrapUserPassword := fmt.Sprintf("this_is_a_very_sensitive_password_%s", hex.EncodeToString(randomBytes[:]))
	bootstrapUserPasswordBcrypt, err := generateBcryptHash(bootstrapUserPassword)
	if err != nil {
		panic(fmt.Errorf("bcrypt bootstrap user password: %w", err))
	}

	// create the bootstrap user inside the dogfood organization
	bootstrapUserEmail := "root@openauth.example.com"
	if _, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:             uuid.New(),
		OrganizationID: dogfoodOrganizationID,
		VerifiedEmail:  &bootstrapUserEmail,
		PasswordBcrypt: &bootstrapUserPasswordBcrypt,
	}); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// create session signing key for the new project
	sessionSigningKey, err := s.CreateSessionSigningKey(ctx, formattedDogfoodProjectID)
	if err != nil {
		return nil, fmt.Errorf("create project signing key: %w", err)
	}

	// create intermediate session signing key for the new project
	intermediateSessionSigningKey, err := s.CreateIntermediateSessionSigningKey(ctx, formattedDogfoodProjectID)
	if err != nil {
		return nil, fmt.Errorf("create project intermediate session signing key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &CreateDogfoodProjectResponse{
		DogfoodProjectID:                   idformat.Project.Format(dogfoodProjectID),
		BootstrapUserEmail:                 bootstrapUserEmail,
		BootstrapUserVerySensitivePassword: bootstrapUserPassword,
		SessionSigningKeyID:                idformat.SessionSigningKey.Format(sessionSigningKey.ID),
		IntermediateSessionSigningKeyID:    idformat.IntermediateSession.Format(intermediateSessionSigningKey.ID),
	}, nil
}
