package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	oauthv1 "github.com/tesseral-labs/tesseral/internal/oauth/gen/tesseral/oauth/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Store) GetSessionPublicKeysByProjectID(ctx context.Context, projectId string) ([]*oauthv1.SessionSigningKey, error) {
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

	var out []*oauthv1.SessionSigningKey
	for _, sessionSigningKey := range sessionSigningKeys {
		pub, err := x509.ParsePKIXPublicKey(sessionSigningKey.PublicKey)
		if err != nil {
			panic(fmt.Errorf("public key from bytes: %w", err))
		}

		jwk, err := structpb.NewStruct(map[string]any{
			"kid": idformat.SessionSigningKey.Format(sessionSigningKey.ID),
			"kty": "EC",
			"crv": "P-256",
			"x":   base64.RawURLEncoding.EncodeToString(pub.(*ecdsa.PublicKey).X.Bytes()),
			"y":   base64.RawURLEncoding.EncodeToString(pub.(*ecdsa.PublicKey).Y.Bytes()),
		})
		if err != nil {
			panic(fmt.Errorf("marshal public key to structpb: %w", err))
		}

		out = append(out, &oauthv1.SessionSigningKey{
			Id:           idformat.SessionSigningKey.Format(sessionSigningKey.ID),
			ProjectId:    idformat.Project.Format(projectID),
			PublicKeyJwk: jwk,
		})
	}

	return out, nil
}
