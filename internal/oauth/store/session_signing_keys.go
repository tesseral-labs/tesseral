package store

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/openauth/openauth/internal/crypto/ecdsa"
	openauthv1 "github.com/openauth/openauth/internal/gen/openauth/v1"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/structpb"
)

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
		pub, err := ecdsa.PublicKeyFromBytes(sessionSigningKey.PublicKey)
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
