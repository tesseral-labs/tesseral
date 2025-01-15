package store

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/common/apierror"
	openauthecdsa "github.com/openauth/openauth/internal/crypto/ecdsa"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Store) GetSessionPublicKeysByProjectID(ctx context.Context, projectId string) ([]*backendv1.SessionSigningKey, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID, err := idformat.Project.Parse(projectId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid project id", fmt.Errorf("parse project id: %w", err))
	}

	sessionSigningKeys, err := q.GetSessionSigningKeysByProjectID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("session signing keys not found", fmt.Errorf("get session signing keys by project id: %w", err))
		}

		return nil, fmt.Errorf("get session signing keys by project id: %w", err)
	}

	var out []*backendv1.SessionSigningKey
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

		out = append(out, &backendv1.SessionSigningKey{
			Id:           idformat.SessionSigningKey.Format(sessionSigningKey.ID),
			PublicKeyJwk: jwk,
		})
	}

	return out, nil
}
