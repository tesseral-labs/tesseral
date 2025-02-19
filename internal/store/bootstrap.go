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

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/bcryptcost"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/store/queries"
	"golang.org/x/crypto/bcrypt"
)

type CreateDogfoodProjectResponse struct {
	DogfoodProjectID                   string
	BootstrapUserEmail                 string
	BootstrapUserVerySensitivePassword string
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
	authDomain := fmt.Sprintf("%s.%s", idformat.Project.Format(dogfoodProjectID), "tesseral.example.app")
	customAuthDomain := "auth.app.tesseral.example.com"

	if _, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                dogfoodProjectID,
		DisplayName:       "OpenAuth Dogfood",
		OrganizationID:    nil, // will populate after creating org
		LogInWithPassword: true,
		LogInWithGoogle:   true,
		AuthDomain:        &authDomain,
		CustomAuthDomain:  &customAuthDomain,
	}); err != nil {
		return nil, fmt.Errorf("create dogfood project: %w", err)
	}

	if _, err := q.CreateProjectUISettings(ctx, queries.CreateProjectUISettingsParams{
		ID:        uuid.New(),
		ProjectID: dogfoodProjectID,
	}); err != nil {
		return nil, fmt.Errorf("create dogfood project ui settings: %w", err)
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
	bootstrapUserPasswordBcryptBytes, err := bcrypt.GenerateFromPassword([]byte(bootstrapUserPassword), bcryptcost.Cost)
	if err != nil {
		panic(fmt.Errorf("bcrypt bootstrap user password: %w", err))
	}

	// create the bootstrap user inside the dogfood organization
	bootstrapUserEmail := "root@tesseral.example.com"
	bootstrapUserPasswordBcrypt := string(bootstrapUserPasswordBcryptBytes)
	if _, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:             uuid.New(),
		OrganizationID: dogfoodOrganizationID,
		Email:          bootstrapUserEmail,
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

	// Encrypt the symmetric key with the KMS
	sskEncryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.intermediateSessionSigningKeyKMSKeyID,
		Plaintext:           privateKeyBytes,
	})
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, err
	}

	// Store the encrypted key in the database
	if _, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            dogfoodProjectID,
		ExpireTime:           &expiresAt,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: sskEncryptOutput.CiphertextBlob,
	}); err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &CreateDogfoodProjectResponse{
		DogfoodProjectID:                   idformat.Project.Format(dogfoodProjectID),
		BootstrapUserEmail:                 bootstrapUserEmail,
		BootstrapUserVerySensitivePassword: bootstrapUserPassword,
	}, nil
}
