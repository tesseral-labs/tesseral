package store

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	openauthecdsa "github.com/openauth/openauth/internal/crypto/ecdsa"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
)

type IntermediateSessionSigningKey struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	CreateTime time.Time
	ExpireTime time.Time
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

func (s *Store) CreateIntermediateSessionSigningKey(ctx context.Context, projectId string) (*IntermediateSessionSigningKey, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(projectId)
	if err != nil {
		return nil, err
	}

	// Allow this key to be used for 7 hours
	// - this adds a 1 hour buffer to the 6 hour key rotation period,
	//   so that the key can be rotated before it expires without
	//   causing existing JWT parsing to fail
	expiresAt := time.Now().Add(time.Hour * 7)

	// Generate a new symmetric key
	privateKey, err := openauthecdsa.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := openauthecdsa.PrivateKeyBytes(privateKey)
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	encryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.intermediateSessionSigningKeyKMSKeyID,
		Plaintext:           privateKeyBytes,
	})
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := openauthecdsa.PublicKeyBytes(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	// Store the encrypted key in the database
	createdIntermediateSessionSigningKey, err := q.CreateIntermediateSessionSigningKey(ctx, queries.CreateIntermediateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            projectID,
		ExpireTime:           &expiresAt,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: encryptOutput.CiphertextBlob,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseIntermediateSessionSigningKey(&createdIntermediateSessionSigningKey, privateKey), nil
}

func parseIntermediateSessionSigningKey(issk *queries.IntermediateSessionSigningKey, privateKey *ecdsa.PrivateKey) *IntermediateSessionSigningKey {
	publicKey := &privateKey.PublicKey

	return &IntermediateSessionSigningKey{
		ID:         issk.ID,
		ProjectID:  issk.ProjectID,
		CreateTime: *issk.CreateTime,
		ExpireTime: *issk.ExpireTime,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}
