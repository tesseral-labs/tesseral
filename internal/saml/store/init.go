package store

import (
	"context"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/saml/authn"
	"github.com/tesseral-labs/tesseral/internal/saml/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
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

	qProject, err := q.GetProject(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	qSAMLConnection, err := q.GetSAMLConnection(ctx, queries.GetSAMLConnectionParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        samlConnectionUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("get saml connection: %w", err)
	}

	spEntityID := fmt.Sprintf("https://%s/api/saml/v1/%s", qProject.VaultDomain, samlConnectionID)

	return &SAMLConnectionInitData{
		SPEntityID:     spEntityID,
		IDPRedirectURL: *qSAMLConnection.IdpRedirectUrl,
	}, nil
}
