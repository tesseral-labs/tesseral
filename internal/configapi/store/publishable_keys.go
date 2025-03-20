package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type GetPublishableKeyConfigurationResponse struct {
	ProjectID      string           `json:"projectId"`
	VaultDomain    string           `json:"vaultDomain"`
	DevMode        bool             `json:"devMode"`
	TrustedDomains []string         `json:"trustedDomains"`
	Keys           []map[string]any `json:"keys"`
}

func (s *Store) GetPublishableKeyConfiguration(ctx context.Context, publishableKey string) (*GetPublishableKeyConfigurationResponse, error) {
	id, err := idformat.PublishableKey.Parse(publishableKey)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid publishable key id", fmt.Errorf("parse publishable key id: %w", err))
	}

	qConfig, err := s.q.GetPublishableKeyConfiguration(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("publishable key not found", fmt.Errorf("get publishable key configuration: %w", err))
		}
		return nil, fmt.Errorf("get publishable key configuration: %w", err)
	}

	trustedDomains, err := s.q.GetPublishableKeyTrustedDomains(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get publishable key trusted domains: %w", err)
	}

	qSessionSigningKeys, err := s.q.GetPublishableKeySessionSigningPublicKeys(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get publishable key session signing public keys: %w", err)
	}

	var keys []map[string]any
	for _, qSessionSigningKey := range qSessionSigningKeys {
		pub, err := x509.ParsePKIXPublicKey(qSessionSigningKey.PublicKey)
		if err != nil {
			panic(fmt.Errorf("public key from bytes: %w", err))
		}

		keys = append(keys, map[string]any{
			"kid": idformat.SessionSigningKey.Format(qSessionSigningKey.ID),
			"kty": "EC",
			"crv": "P-256",
			"x":   base64.RawURLEncoding.EncodeToString(pub.(*ecdsa.PublicKey).X.Bytes()),
			"y":   base64.RawURLEncoding.EncodeToString(pub.(*ecdsa.PublicKey).Y.Bytes()),
		})
	}

	return &GetPublishableKeyConfigurationResponse{
		ProjectID:      idformat.Project.Format(qConfig.ProjectID),
		VaultDomain:    qConfig.VaultDomain,
		DevMode:        qConfig.DevMode,
		TrustedDomains: trustedDomains,
		Keys:           keys,
	}, nil
}
