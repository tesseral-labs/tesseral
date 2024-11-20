package store

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/google/uuid"
	openauthecdsa "github.com/openauth-dev/openauth/internal/ecdsa"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type IntermediateSessionSigningKey struct {
	ID 							uuid.UUID
	ProjectID 			uuid.UUID
	CreateTime 			time.Time
	ExpireTime 			time.Time
	PublicKey 			*ecdsa.PublicKey
	PrivateKey 			*ecdsa.PrivateKey
}

func (s *Store) CreateIntermediateSessionSigningKey(ctx context.Context, projectID string) (*queries.IntermediateSessionSigningKey, error) {
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
	ecdsaKeyPair, err := openauthecdsa.New()
	if err != nil {
		return nil, err
	}

	// Extract the private key
	privateKeyBytes, err := ecdsaKeyPair.PrivateKeyPEM()
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	encrytpOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		KeyId:    	&s.intermediateSessionSigningKeyKMSKeyID,
		Plaintext: 	privateKeyBytes,
	})
	if err != nil {
		return nil, err
	}

	publicKey, err := ecdsaKeyPair.PublicKeyPEM()
	if err != nil {
		return nil, err
	}

	// Store the encrypted key in the database
	createdIntermediateSessionSigningKey, err := q.CreateIntermediateSessionSigningKey(ctx, queries.CreateIntermediateSessionSigningKeyParams{
		ID: uuid.New(),
		ProjectID: projectId,
		ExpireTime: &expiresAt,
		PublicKey: publicKey,
		PrivateKeyCipherText: encrytpOutput.CipherTextBlob,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &createdIntermediateSessionSigningKey, nil
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
	signingKey, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: intermediateSessionSigningKey.PrivateKeyCipherText,
		KeyId: &s.intermediateSessionSigningKeyKMSKeyID,
	})
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted signing key
	ecdsaKeyPair, err := openauthecdsa.NewFromBytes(signingKey.Value)
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted private key
	return &IntermediateSessionSigningKey{
		ID: intermediateSessionSigningKey.ID,
		ProjectID: intermediateSessionSigningKey.ProjectID,
		CreateTime: *intermediateSessionSigningKey.CreateTime,
		ExpireTime: *intermediateSessionSigningKey.ExpireTime,
		PublicKey: ecdsaKeyPair.PublicKey,
		PrivateKey: ecdsaKeyPair.PrivateKey,
	}, nil
}
