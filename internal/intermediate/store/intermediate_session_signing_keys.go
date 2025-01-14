package store

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	openauthecdsa "github.com/openauth/openauth/internal/crypto/ecdsa"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

type IntermediateSessionSigningKey struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	CreateTime time.Time
	ExpireTime time.Time
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

func (s *Store) GetIntermediateSessionSigningKeyByID(ctx context.Context, id string) (*IntermediateSessionSigningKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

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
		CiphertextBlob:      intermediateSessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.intermediateSessionSigningKeyKMSKeyID,
	})
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted signing key
	privateKey, err := openauthecdsa.PrivateKeyFromBytes(decryptOutput.Plaintext)
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted private key
	return parseIntermediateSessionSigningKey(&intermediateSessionSigningKey, privateKey), nil
}

func (s *Store) GetIntermediateSessionSigningKeyByProjectID(ctx context.Context, projectId string) (*IntermediateSessionSigningKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(projectId)
	if err != nil {
		return nil, err
	}

	// Fetch the raw record from the database
	intermediateSessionSigningKey, err := q.GetIntermediateSessionSigningKeyByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Decrypt the signing key using KMS
	decryptOutput, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      intermediateSessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.intermediateSessionSigningKeyKMSKeyID,
	})
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted signing key
	privateKey, err := openauthecdsa.PrivateKeyFromBytes(decryptOutput.Plaintext)
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted private key
	return parseIntermediateSessionSigningKey(&intermediateSessionSigningKey, privateKey), nil
}

func (s *Store) GetIntermediateSessionPublicKeyByProjectID(ctx context.Context, projectId string) (*ecdsa.PublicKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(projectId)
	if err != nil {
		return nil, err
	}

	// Fetch the raw record from the database
	intermediateSessionSigningKey, err := q.GetIntermediateSessionSigningKeyByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted signing key
	publicKey, err := openauthecdsa.PublicKeyFromBytes(intermediateSessionSigningKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
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
