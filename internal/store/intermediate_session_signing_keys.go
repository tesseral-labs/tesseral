package store

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	openauthecdsa "github.com/openauth-dev/openauth/internal/crypto/ecdsa"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type IntermediateSessionSigningKey struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	CreateTime time.Time
	ExpireTime time.Time
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

func (s *Store) CreateIntermediateSessionSigningKey(ctx context.Context, projectID string) (*IntermediateSessionSigningKey, error) {
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
		KeyId:     &s.intermediateSessionSigningKeyKMSKeyID,
		Plaintext: privateKeyBytes,
	})
	if err != nil {
		return nil, err
	}

	// Store the encrypted key in the database
	createdIntermediateSessionSigningKey, err := q.CreateIntermediateSessionSigningKey(ctx, queries.CreateIntermediateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            projectId,
		ExpireTime:           &expiresAt,
		PrivateKeyCipherText: encryptOutput.CipherTextBlob,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseIntermediateSessionSigningKey(&createdIntermediateSessionSigningKey, privateKey), nil
}

func (s *Store) GetIntermediateSessionSigningKeyByID(ctx context.Context, id string) (*IntermediateSessionSigningKey, error) {
	_, q, _, _, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}

	intermediateSessionSigningKeyID, err := idformat.IntermediateSessionSigningKey.Parse(id)
	if err != nil {
		return nil, err
	}

	// Fetch the raw record from the database
	intermediateSessionSigningKey, err := q.GetIntermediateSessionSigningKeyByID(ctx, intermediateSessionSigningKeyID)
	if err != nil {
		return nil, err
	}

	// Decrypt the signing key using KMS
	decryptOutput, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: intermediateSessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:          &s.intermediateSessionSigningKeyKMSKeyID,
	})
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted signing key
	privateKey, err := openauthecdsa.PrivateKeyFromBytes(decryptOutput.Value)
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted private key
	return parseIntermediateSessionSigningKey(&intermediateSessionSigningKey, privateKey), nil
}

func (s *Store) GetIntermediateSessionSigningKeyByProjectID(ctx context.Context, projectIdString string) (*IntermediateSessionSigningKey, error) {
	_, q, _, _, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}

	projectId, err := idformat.Project.Parse(projectIdString)
	if err != nil {
		return nil, err
	}

	// Fetch the raw record from the database
	intermediateSessionSigningKey, err := q.GetIntermediateSessionSigningKeyByProjectID(ctx, projectId)
	if err != nil {
		return nil, err
	}

	// Decrypt the signing key using KMS
	decryptOutput, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: intermediateSessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:          &s.intermediateSessionSigningKeyKMSKeyID,
	})
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted signing key
	privateKey, err := openauthecdsa.PrivateKeyFromBytes(decryptOutput.Value)
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted private key
	return parseIntermediateSessionSigningKey(&intermediateSessionSigningKey, privateKey), nil
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
