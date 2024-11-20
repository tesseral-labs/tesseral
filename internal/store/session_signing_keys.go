package store

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/rsakeys"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type SessionSigningKey struct {
	ID 						uuid.UUID
	ProjectID 		uuid.UUID
	CreateTime 		time.Time
	ExpireTime 		time.Time
	PublicKey 		[]byte
	PrivateKey 		[]byte
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

	// Allow 15 minutes for the user to verify their email before expiring the intermediate session
	expiresAt := time.Now().Add(time.Minute * 15)

	// Generate a new symmetric key
	privateKey, publicKey, err := rsakeys.GenerateRSAKeys()
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	keyId := "test" // TODO: get the key ID from the environment or context
	encrytpOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		KeyId:    	&keyId,
		Plaintext: 	privateKey,
	})
	if err != nil {
		return nil, err
	}


	// Create the new method verification challenge
	sessionSigningKey, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID: uuid.New(),
		ProjectID: projectId,
		ExpireTime: &expiresAt,
		PublicKey: publicKey,
		PrivateKeyCipherText: encrytpOutput.CipherTextBlob,
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
	})
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted signing key
	return &SessionSigningKey{
		ID: sessionSigningKey.ID,
		ProjectID: sessionSigningKey.ProjectID,
		CreateTime: *sessionSigningKey.CreateTime,
		ExpireTime: *sessionSigningKey.ExpireTime,
		PublicKey: sessionSigningKey.PublicKey,
		PrivateKey: signingKey.Value,
	}, nil
}