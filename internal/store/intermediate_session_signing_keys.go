package store

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
	"github.com/openauth-dev/openauth/internal/symmetrickeys"
)

type IntermediateSessionSigningKey struct {
	ID string
	ProjectID string
	CreateTime time.Time
	ExpireTime time.Time
	SigningKey string
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
	key, err := symmetrickeys.GenerateSymmetricKey()
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	keyId := "test" // TODO: get the key ID from the environment or context
	encrytpOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		KeyId:    	&keyId,
		Plaintext: 	[]byte(key),
	})
	if err != nil {
		return nil, err
	}

	// Store the encrypted key in the database
	createdIntermediateSessionSigningKey, err := q.CreateIntermediateSessionSigningKey(ctx, queries.CreateIntermediateSessionSigningKeyParams{
		ID: uuid.New(),
		ProjectID: projectId,
		ExpireTime: &expiresAt,
		SigningKeyCipherText: encrytpOutput.CipherTextBlob,
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

	sessionID, err := idformat.IntermediateSessionSigningKey.Parse(id)
	if err != nil {
		return nil, err
	}

	// Fetch the raw record from the database
	intermediateSessionSigningKey, err := q.GetIntermediateSessionSigningKeyByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Decrypt the signing key using KMS
	signingKey, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: intermediateSessionSigningKey.SigningKeyCipherText,
	})
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted signing key
	return &IntermediateSessionSigningKey{
		ID: intermediateSessionSigningKey.ID.String(),
		ProjectID: intermediateSessionSigningKey.ProjectID.String(),
		CreateTime: *intermediateSessionSigningKey.CreateTime,
		ExpireTime: *intermediateSessionSigningKey.ExpireTime,
		SigningKey: string(signingKey.PlainText),
	}, nil
}
