package store

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/saml/authn"
	"github.com/tesseral-labs/tesseral/internal/saml/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

const sessionDuration = time.Hour * 24 * 7

type SAMLConnectionACSData struct {
	IDPX509Certificate  *x509.Certificate
	IDPEntityID         string
	SPEntityID          string
	OrganizationID      string
	OrganizationDomains []string
}

func (s *Store) GetSAMLConnectionACSData(ctx context.Context, samlConnectionID string) (*SAMLConnectionACSData, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	samlConnectionUUID, err := idformat.SAMLConnection.Parse(samlConnectionID)
	if err != nil {
		return nil, fmt.Errorf("parse saml connection id: %w", err)
	}

	qSAMLConnection, err := q.GetSAMLConnection(ctx, queries.GetSAMLConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        samlConnectionUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("get saml connection: %w", err)
	}

	idpX509Certificate, err := x509.ParseCertificate(qSAMLConnection.IdpX509Certificate)
	if err != nil {
		panic(fmt.Errorf("parse idp x509 certificate: %w", err))
	}

	spEntityID := fmt.Sprintf("https://%s/api/saml/v1/%s", qProject.VaultDomain, samlConnectionID)

	organizationDomains, err := q.GetOrganizationDomains(ctx, qSAMLConnection.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get organization domains: %w", err)
	}

	return &SAMLConnectionACSData{
		IDPX509Certificate:  idpX509Certificate,
		IDPEntityID:         *qSAMLConnection.IdpEntityID,
		SPEntityID:          spEntityID,
		OrganizationDomains: organizationDomains,
		OrganizationID:      idformat.Organization.Format(qSAMLConnection.OrganizationID),
	}, nil
}

type CreateSessionRequest struct {
	SAMLConnectionID string
	Email            string
}

type CreateSessionResponse struct {
	ProjectCookieDomain string
	RedirectURI         string
	RefreshToken        string
}

func (s *Store) CreateSession(ctx context.Context, req *CreateSessionRequest) (*CreateSessionResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	samlConnectionID, err := idformat.SAMLConnection.Parse(req.SAMLConnectionID)
	if err != nil {
		return nil, fmt.Errorf("parse saml connection id: %w", err)
	}

	qSAMLConnection, err := q.GetSAMLConnection(ctx, queries.GetSAMLConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        samlConnectionID,
	})
	if err != nil {
		return nil, fmt.Errorf("get saml connection: %w", err)
	}

	qUser, err := s.upsertUser(ctx, q, qSAMLConnection.OrganizationID, req.Email)
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
			SamlConnectionId: refOrNil(idformat.SAMLConnection.Format(samlConnectionID)),
		},
		ResourceType:   queries.AuditLogEventResourceTypeSession,
		ResourceID:     &qSession.ID,
		OrganizationID: &qSAMLConnection.OrganizationID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	redirectURI := qProject.RedirectUri
	if qProject.AfterLoginRedirectUri != nil {
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
