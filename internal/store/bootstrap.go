package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/crypto/ecdsa"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
)

type CreateDogfoodProjectResponse struct {
	DogfoodProjectID                   string
	BootstrapUserEmail                 string
	BootstrapUserVerySensitivePassword string
	SessionSigningKeyID                string
	IntermediateSessionSigningKeyID    string
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

	if _, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                       dogfoodProjectID,
		DisplayName:              "OpenAuth Dogfood",
		OrganizationID:           nil, // will populate after creating org
		LogInWithPasswordEnabled: true,
		LogInWithGoogleEnabled:   true,
	}); err != nil {
		return nil, fmt.Errorf("create dogfood project: %w", err)
	}

	if _, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                   dogfoodOrganizationID,
		ProjectID:            dogfoodProjectID,
		DisplayName:          "OpenAuth",
		OverrideLogInMethods: false, // just to be explicit here
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
		Email:          bootstrapUserEmail,
		IsOwner:        true,
		PasswordBcrypt: &bootstrapUserPasswordBcrypt,
	}); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// create session signing keys for the new project
	// Allow this key to be used for 7 hours
	// - this adds a 1 hour buffer to the 6 hour key rotation period,
	//   so that the key can be rotated before it expires without
	//   causing existing JWT parsing to fail
	expiresAt := time.Now().Add(time.Hour * 7)

	// Generate a new symmetric key
	privateKey, err := ecdsa.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := ecdsa.PrivateKeyBytes(privateKey)
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

	publicKeyBytes, err := ecdsa.PublicKeyBytes(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	isskEncryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
		Plaintext:           privateKeyBytes,
	})
	if err != nil {
		return nil, err
	}

	// Store the encrypted key in the database
	sessionSigningKey, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            dogfoodProjectID,
		ExpireTime:           &expiresAt,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: sskEncryptOutput.CipherTextBlob,
	})
	if err != nil {
		return nil, err
	}

	intermediateSessionSigningKey, err := q.CreateIntermediateSessionSigningKey(ctx, queries.CreateIntermediateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            dogfoodProjectID,
		ExpireTime:           &expiresAt,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: isskEncryptOutput.CipherTextBlob,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &CreateDogfoodProjectResponse{
		DogfoodProjectID:                   idformat.Project.Format(dogfoodProjectID),
		BootstrapUserEmail:                 bootstrapUserEmail,
		BootstrapUserVerySensitivePassword: bootstrapUserPassword,
		SessionSigningKeyID:                idformat.SessionSigningKey.Format(sessionSigningKey.ID),
		IntermediateSessionSigningKeyID:    idformat.IntermediateSessionSigningKey.Format(intermediateSessionSigningKey.ID),
	}, nil
}
