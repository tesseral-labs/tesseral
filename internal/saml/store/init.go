package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/saml/authn"
	"github.com/tesseral-labs/tesseral/internal/saml/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type SAMLConnectionInitData struct {
	SPEntityID     string
	IDPRedirectURL string
}

func (s *Store) GetSAMLConnectionInitData(ctx context.Context, samlConnectionID string) (*SAMLConnectionInitData, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
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

	auditSamlConnection, err := s.auditlogStore.GetSAMLConnection(ctx, tx, samlConnectionUUID)
	if err != nil {
		return nil, fmt.Errorf("get audit saml connection: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.saml_connections.initiate",
		EventDetails: &auditlogv1.InitiateSAMLConnection{
			SamlConnection: auditSamlConnection,
		},
		ResourceType:   queries.AuditLogEventResourceTypeSamlConnection,
		ResourceID:     (*uuid.UUID)(&samlConnectionUUID),
		OrganizationID: &qSAMLConnection.OrganizationID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	spEntityID := fmt.Sprintf("https://%s/api/saml/v1/%s", qProject.VaultDomain, samlConnectionID)

	return &SAMLConnectionInitData{
		SPEntityID:     spEntityID,
		IDPRedirectURL: *qSAMLConnection.IdpRedirectUrl,
	}, nil
}
