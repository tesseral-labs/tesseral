package store

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	openauthecdsa "github.com/openauth/openauth/internal/crypto/ecdsa"
	openauthv1 "github.com/openauth/openauth/internal/gen/openauth/v1"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
	"google.golang.org/protobuf/types/known/structpb"
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
	privateKey, err := openauthecdsa.GenerateKey()
	if err != nil {
		return nil, err
	}

	// Extract the private key
	privateKeyBytes, err := openauthecdsa.PrivateKeyBytes(privateKey)
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

	publicKeyBytes, err := openauthecdsa.PublicKeyBytes(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	// Create the new method verification challenge
	sessionSigningKey, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            projectId,
		ExpireTime:           &expiresAt,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: encryptOutput.CipherTextBlob,
	})
	if err != nil {
		return nil, err
	}

	slog.Info("sessionSigningKey", "sessionSigningKey", sessionSigningKey)

	// Commit the transaction
	if err := commit(); err != nil {
		return nil, err
	}

	// Return the new method verification challenge
	return parseSessionSigningKey(&sessionSigningKey, privateKey), nil
}

func (s *Store) GetSessionSigningKeyByID(ctx context.Context, id string) (*SessionSigningKey, error) {
	_, q, _, _, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}

	sessionSigningKeyID, err := idformat.SessionSigningKey.Parse(id)
	if err != nil {
		return nil, err
	}

	// Fetch the raw record from the database
	sessionSigningKey, err := q.GetSessionSigningKeyByID(ctx, sessionSigningKeyID)
	if err != nil {
		return nil, err
	}

	// Decrypt the signing key using KMS
	decryptOutput, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      sessionSigningKey.PrivateKeyCipherText,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
	})
	if err != nil {
		return nil, err
	}

	// Create an ECDSA key pair from the decrypted private key
	privateKey, err := openauthecdsa.PrivateKeyFromBytes(decryptOutput.Value)
	if err != nil {
		return nil, err
	}

	// Return the intermediate session signing key with the decrypted signing key
	return parseSessionSigningKey(&sessionSigningKey, privateKey), nil
}

func (s *Store) GetSessionPublicKeysByProjectID(ctx context.Context, projectId string) ([]*openauthv1.SessionSigningKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(projectId)
	if err != nil {
		return nil, err
	}

	sessionSigningKeys, err := q.GetSessionSigningKeysByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var out []*openauthv1.SessionSigningKey
	for _, sessionSigningKey := range sessionSigningKeys {
		pub, err := openauthecdsa.PublicKeyFromBytes(sessionSigningKey.PublicKey)
		if err != nil {
			panic(fmt.Errorf("public key from bytes: %w", err))
		}

		jwk, err := structpb.NewStruct(map[string]any{
			"kid": idformat.SessionSigningKey.Format(sessionSigningKey.ID),
			"kty": "EC",
			"crv": "P-256",
			"x":   base64.RawURLEncoding.EncodeToString(pub.X.Bytes()),
			"y":   base64.RawURLEncoding.EncodeToString(pub.Y.Bytes()),
		})
		if err != nil {
			panic(fmt.Errorf("marshal public key to structpb: %w", err))
		}

		out = append(out, &openauthv1.SessionSigningKey{
			Id:           idformat.SessionSigningKey.Format(sessionSigningKey.ID),
			ProjectId:    idformat.Project.Format(projectID),
			PublicKeyJwk: jwk,
		})
	}

	return out, nil
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
