package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/saml/authn"
	"github.com/openauth/openauth/internal/saml/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

type SAMLConnectionInitData struct {
	SPEntityID     string
	IDPRedirectURL string
}

func (s *Store) GetSAMLConnectionInitData(ctx context.Context, samlConnectionID string) (*SAMLConnectionInitData, error) {
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

	spEntityID := fmt.Sprintf("http://localhost:3001/saml/v1/%s", samlConnectionID) // todo

	return &SAMLConnectionInitData{
		SPEntityID:     spEntityID,
		IDPRedirectURL: *qSAMLConnection.IdpRedirectUrl,
	}, nil
}
