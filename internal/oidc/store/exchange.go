package store

import (
	"context"
	"encoding/base64"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/oidc/authn"
	"github.com/tesseral-labs/tesseral/internal/oidc/store/queries"
	"github.com/tesseral-labs/tesseral/internal/oidcclient"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type OIDCUserData struct {
	Email               string
	OrganizationID      string
	OrganizationDomains []string
	RedirectURL         string
}

func (s *Store) ExchangeOIDCCode(ctx context.Context, oidcConnectionID string, code string) (*OIDCUserData, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	oidcConnectionUUID, err := idformat.OIDCConnection.Parse(oidcConnectionID)
	if err != nil {
		return nil, fmt.Errorf("parse oidc connection id: %w", err)
	}

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	qOIDCConnection, err := q.GetOIDCConnection(ctx, queries.GetOIDCConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        oidcConnectionUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	organizationDomains, err := q.GetOrganizationDomains(ctx, qOIDCConnection.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get organization domains: %w", err)
	}

	config, err := s.oidc.GetConfiguration(ctx, qOIDCConnection.ConfigurationUrl)
	if err != nil {
		return nil, fmt.Errorf("get OIDC configuration: %w", err)
	}

	var (
		clientAuthBasic string
		clientAuthPost  string
	)
	if qOIDCConnection.ClientSecretCiphertext != nil {
		decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
			KeyId:               &s.oidcClientSecretsKMSKeyID,
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
			CiphertextBlob:      qOIDCConnection.ClientSecretCiphertext,
		})
		if err != nil {
			return nil, fmt.Errorf("decrypt oidc client secret: %w", err)
		}
		switch {
		case slices.Contains(config.TokenEndpointAuthMethodsSupported, "client_secret_post"):
			clientAuthPost = string(decryptRes.Plaintext)
		case slices.Contains(config.TokenEndpointAuthMethodsSupported, "client_secret_basic") || len(config.TokenEndpointAuthMethodsSupported) == 0: // If omitted, the default is client_secret_basic
			clientAuthBasic = base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", qOIDCConnection.ClientID, string(decryptRes.Plaintext)))
		default:
			return nil, fmt.Errorf("OIDC connection %s does not support client authentication method for token endpoint: %s", oidcConnectionID, config.TokenEndpointAuthMethodsSupported)
		}
	}

	tokenRes, err := s.oidc.ExchangeCode(ctx, oidcclient.ExchangeCodeRequest{
		TokenEndpoint:   config.TokenEndpoint,
		Code:            code,
		RedirectURI:     fmt.Sprintf("https://%s/api/oidc/v1/%s/callback", qProject.VaultDomain, idformat.OIDCConnection.Format(qOIDCConnection.ID)),
		ClientID:        qOIDCConnection.ClientID,
		ClientAuthBasic: clientAuthBasic,
		ClientAuthPost:  clientAuthPost,
		CodeVerifier:    authn.IntermediateSession(ctx).OidcCodeVerifier,
	})
	if err != nil {
		return nil, fmt.Errorf("exchange OIDC code: %w", err)
	}

	claims, err := s.oidc.ValidateIDToken(ctx, oidcclient.ValidateIDTokenRequest{
		Configuration: config,
		IDToken:       tokenRes.IDToken,
	})
	if err != nil {
		return nil, fmt.Errorf("validate id token: %w", err)
	}

	if err := q.UpdateIntermediateSession(ctx, queries.UpdateIntermediateSessionParams{
		ID:                       authn.IntermediateSession(ctx).ID,
		Email:                    &claims.Email,
		VerifiedOidcConnectionID: (*uuid.UUID)(&oidcConnectionUUID),
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &OIDCUserData{
		OrganizationID:      idformat.Organization.Format(qOIDCConnection.OrganizationID),
		OrganizationDomains: organizationDomains,
		Email:               claims.Email,
		RedirectURL:         fmt.Sprintf("https://%s/finish-login", qProject.VaultDomain),
	}, nil
}
