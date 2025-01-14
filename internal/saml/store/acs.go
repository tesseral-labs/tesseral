package store

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/openauth/openauth/internal/saml/authn"
	"github.com/openauth/openauth/internal/saml/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
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

	spEntityID := fmt.Sprintf("http://localhost:3001/saml/v1/%s", samlConnectionID) // todo

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
