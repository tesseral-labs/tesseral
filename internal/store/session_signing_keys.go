package store

import (
	"context"
	"crypto/ecdsa"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/google/uuid"
	openauthecdsa "github.com/openauth-dev/openauth/internal/crypto/ecdsa"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type SessionSigningKey struct {
	ID 						uuid.UUID
	ProjectID 		uuid.UUID
	CreateTime 		time.Time
	ExpireTime 		time.Time
	PrivateKey 		*ecdsa.PrivateKey
	PublicKey 		*ecdsa.PublicKey
}

func (s *Store) CreateSessionSigningKey(ctx context.Context, projectID string) (*queries.SessionSigningKey, error) {
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
	privateKey, err := ecdsaKeyPair.PrivateKeyPEM()
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	encryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		KeyId:    	&s.sessionSigningKeyKmsKeyID,
		Plaintext: 	privateKey,
	})
	if err != nil {
		return nil, err
	}

	// Extract the public key
	publicKey, err := ecdsaKeyPair.PublicKeyPEM()
	if err != nil {
		return nil, err
	}

	// Create the new method verification challenge
	sessionSigningKey, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID: uuid.New(),
		ProjectID: projectId,
		ExpireTime: &expiresAt,
		PublicKey: publicKey,
		PrivateKeyCipherText: encryptOutput.CipherTextBlob,
	})
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := commit(); err != nil {
		return nil, err
	}

	// Return the new method verification challenge
	return &sessionSigningKey, nil
}

func (s *Store) GetSessionSigningKeyByID(ctx context.Context, id string) (*SessionSigningKey, error) {
	_, q, _, _, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}

	sessionSigningKeyID, err := idformat.IntermediateSessionSigningKey.Parse(id)
	if err != nil {
		return nil, err
	}

	// Fetch the raw record from the database
	sessionSigningKey, err := q.GetIntermediateSessionSigningKeyByID(ctx, sessionSigningKeyID)
	if err != nil {
		return nil, err
	}

	// Decrypt the signing key using KMS
	signingKey, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: sessionSigningKey.PrivateKeyCipherText,
		KeyId: &s.sessionSigningKeyKmsKeyID,
	})
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted private key
	ecdsaKeyPair, err := openauthecdsa.NewFromBytes(signingKey.Value)
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted signing key
	return &SessionSigningKey{
		ID: sessionSigningKey.ID,
		ProjectID: sessionSigningKey.ProjectID,
		CreateTime: *sessionSigningKey.CreateTime,
		ExpireTime: *sessionSigningKey.ExpireTime,
		PrivateKey: ecdsaKeyPair.PrivateKey,
		PublicKey: ecdsaKeyPair.PublicKey,
	}, nil
}