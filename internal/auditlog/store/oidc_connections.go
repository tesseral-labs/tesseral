package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetOIDCConnection(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.OIDCConnection, error) {
	qOIDCConnection, err := queries.New(db).GetOIDCConnection(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	qOrg, err := queries.New(db).GetOrganization(ctx, qOIDCConnection.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	qProject, err := queries.New(db).GetProject(ctx, qOrg.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	return &auditlogv1.OIDCConnection{
		Id:               idformat.OIDCConnection.Format(qOIDCConnection.ID),
		CreateTime:       timestamppb.New(*qOIDCConnection.CreateTime),
		UpdateTime:       timestamppb.New(*qOIDCConnection.UpdateTime),
		Primary:          &qOIDCConnection.IsPrimary,
		ConfigurationUrl: qOIDCConnection.ConfigurationUrl,
		ClientId:         qOIDCConnection.ClientID,
		RedirectUri:      fmt.Sprintf("https://%s/api/oidc/v1/%s/callback", qProject.VaultDomain, idformat.OIDCConnection.Format(qOIDCConnection.ID)),
	}, nil
}
