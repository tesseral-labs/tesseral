package store

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/oidc/authn"
	"github.com/tesseral-labs/tesseral/internal/oidc/store/queries"
	"github.com/tesseral-labs/tesseral/internal/oidcclient"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

const sessionDuration = time.Hour * 24 * 7

type OIDCUserData struct {
	Email               string
	OrganizationID      string
	OrganizationDomains []string
}

func (s *Store) ExchangeOIDCCode(ctx context.Context, oidcConnectionID string, oidcIntermediateSessionID string, code string) (*OIDCUserData, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	oidcIntermediateSessionUUID, err := idformat.OIDCIntermediateSession.Parse(oidcIntermediateSessionID)
	if err != nil {
		return nil, fmt.Errorf("parse oidc session id: %w", err)
	}

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	qOIDCSession, err := q.DeleteOIDCIntermediateSession(ctx, oidcIntermediateSessionUUID)
	if err != nil {
		return nil, fmt.Errorf("get oidc session: %w", err)
	}
	if idformat.OIDCConnection.Format(qOIDCSession.OidcConnectionID) != oidcConnectionID {
		return nil, fmt.Errorf("oidc intermediate session %s does not match oidc connection %s", oidcIntermediateSessionID, oidcConnectionID)
	}

	qOIDCConnection, err := q.GetOIDCConnection(ctx, queries.GetOIDCConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        qOIDCSession.OidcConnectionID,
	})
	if err != nil {
		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	organizationDomains, err := q.GetOrganizationDomains(ctx, qOIDCConnection.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get organization domains: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
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
		CodeVerifier:    qOIDCSession.CodeVerifier,
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

	return &OIDCUserData{
		OrganizationID:      idformat.Organization.Format(qOIDCConnection.OrganizationID),
		OrganizationDomains: organizationDomains,
		Email:               claims.Email,
	}, nil
}

type CreateSessionRequest struct {
	OIDCConnectionID string
	Email            string
}

type CreateSessionResponse struct {
	ProjectCookieDomain string
	RedirectURI         string
	RefreshToken        string
}

func (s *Store) CreateSession(ctx context.Context, req CreateSessionRequest) (*CreateSessionResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	oidcConnectionID, err := idformat.OIDCConnection.Parse(req.OIDCConnectionID)
	if err != nil {
		return nil, fmt.Errorf("parse oidc connection id: %w", err)
	}

	qOIDCConnection, err := q.GetOIDCConnection(ctx, queries.GetOIDCConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        oidcConnectionID,
	})
	if err != nil {
		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	qUser, err := s.upsertUser(ctx, q, qOIDCConnection.OrganizationID, req.Email)
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}

	expireTime := time.Now().Add(sessionDuration)

	refreshToken := uuid.New()
	refreshTokenSHA256 := sha256.Sum256(refreshToken[:])
	qSession, err := q.CreateSession(ctx, queries.CreateSessionParams{
		ID:                 uuid.Must(uuid.NewV7()),
		ExpireTime:         &expireTime,
		RefreshTokenSha256: refreshTokenSHA256[:],
		UserID:             qUser.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	auditSession, err := s.auditlogStore.GetSession(ctx, tx, qSession.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit session: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.sessions.create",
		EventDetails: &auditlogv1.CreateSession{
			Session:          auditSession,
			OidcConnectionId: &req.OIDCConnectionID,
		},
		ResourceType:   queries.AuditLogEventResourceTypeSession,
		ResourceID:     &qSession.ID,
		OrganizationID: &qOIDCConnection.OrganizationID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	redirectURI := qProject.RedirectUri
	if qProject.AfterLoginRedirectUri != nil && *qProject.AfterLoginRedirectUri != "" {
		redirectURI = *qProject.AfterLoginRedirectUri
	}

	return &CreateSessionResponse{
		ProjectCookieDomain: qProject.CookieDomain,
		RedirectURI:         redirectURI,
		RefreshToken:        idformat.SessionRefreshToken.Format(refreshToken),
	}, nil
}

func (s *Store) upsertUser(ctx context.Context, q *queries.Queries, organizationID uuid.UUID, email string) (*queries.User, error) {
	qUser, err := q.GetUserByEmail(ctx, queries.GetUserByEmailParams{
		OrganizationID: organizationID,
		Email:          email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// upsert a new user instead
			qUser, err := q.CreateUser(ctx, queries.CreateUserParams{
				ID:             uuid.New(),
				OrganizationID: organizationID,
				Email:          email,
				IsOwner:        false,
			})
			if err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}

			return &qUser, nil
		}

		return nil, err
	}

	return &qUser, nil
}
