package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/store/queries"
)

type SessionSigningKey struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	CreateTime time.Time
	ExpireTime time.Time
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func (s *Store) CreateSessionSigningKey(ctx context.Context, projectID string) (*SessionSigningKey, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId, err := idformat.Project.Parse(projectID)
	if err != nil {
		return nil, err
	}

	// Allow this key to be used for 7 hours
	// - this adds a 1 hour buffer to the 6 hour key rotation period,
	//   so that the key can be rotated before it expires without
	//   causing existing JWT parsing to fail
	expiresAt := time.Now().Add(time.Hour * 7)

	// Generate a new symmetric key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Extract the private key
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	encryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
		Plaintext:           privateKeyBytes,
	})
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := commit(); err != nil {
		return nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, err
	}

	sessionSigningKey, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            projectId,
		ExpireTime:           &expiresAt,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: encryptOutput.CiphertextBlob,
	})
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := commit(); err != nil {
		return nil, err
	}

	// Return the new method verification challenge
	return parseSessionSigningKey(&sessionSigningKey, privateKey), nil
}

func parseSessionSigningKey(ssk *queries.SessionSigningKey, privateKey *ecdsa.PrivateKey) *SessionSigningKey {
	publicKey := &privateKey.PublicKey

	return &SessionSigningKey{
		ID:         ssk.ID,
		ProjectID:  ssk.ProjectID,
		CreateTime: *ssk.CreateTime,
		ExpireTime: *ssk.ExpireTime,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}
