package store

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/saml/authn"
	"github.com/tesseral-labs/tesseral/internal/saml/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

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

type FinishLoginRequest struct {
	VerifiedSAMLConnectionID string
	Email                    string
}

func (s *Store) FinishLogin(ctx context.Context, req FinishLoginRequest) (string, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return "", err
	}
	defer rollback()

	samlConnectionUUID, err := idformat.SAMLConnection.Parse(req.VerifiedSAMLConnectionID)
	if err != nil {
		return "", fmt.Errorf("parse saml connection id: %w", err)
	}

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return "", fmt.Errorf("get project: %w", err)
	}

	qSAMLConnection, err := q.GetSAMLConnection(ctx, queries.GetSAMLConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        samlConnectionUUID,
	})
	if err != nil {
		return "", fmt.Errorf("get saml connection: %w", err)
	}

	if err := q.UpdateIntermediateSession(ctx, queries.UpdateIntermediateSessionParams{
		ID:                       *authn.IntermediateSessionID(ctx),
		VerifiedSamlConnectionID: (*uuid.UUID)(&samlConnectionUUID),
		OrganizationID:           &qSAMLConnection.OrganizationID,
		Email:                    &req.Email,
	}); err != nil {
		return "", fmt.Errorf("init intermediate session: %w", err)
	}

	if err := commit(); err != nil {
		return "", fmt.Errorf("commit transaction: %w", err)
	}

	return fmt.Sprintf("https://%s/finish-login", qProject.VaultDomain), nil
}
