package store

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetSAMLConnection(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.SAMLConnection, error) {
	qSAMLConnection, err := queries.New(db).GetSAMLConnection(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get saml connection: %w", err)
	}

	qOrg, err := queries.New(db).GetOrganization(ctx, qSAMLConnection.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	qProject, err := queries.New(db).GetProject(ctx, qOrg.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	var certPEM string
	if len(qSAMLConnection.IdpX509Certificate) != 0 {
		cert, err := x509.ParseCertificate(qSAMLConnection.IdpX509Certificate)
		if err != nil {
			panic(err)
		}

		certPEM = string(pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}))
	}

	spACSURL := fmt.Sprintf("https://%s/api/saml/v1/%s/acs", qProject.VaultDomain, idformat.SAMLConnection.Format(qSAMLConnection.ID))
	spEntityID := fmt.Sprintf("https://%s/api/saml/v1/%s", qProject.VaultDomain, idformat.SAMLConnection.Format(qSAMLConnection.ID))

	return &auditlogv1.SAMLConnection{
		Id:                 idformat.SAMLConnection.Format(qSAMLConnection.ID),
		CreateTime:         timestamppb.New(*qSAMLConnection.CreateTime),
		UpdateTime:         timestamppb.New(*qSAMLConnection.UpdateTime),
		Primary:            &qSAMLConnection.IsPrimary,
		SpAcsUrl:           spACSURL,
		SpEntityId:         spEntityID,
		IdpRedirectUrl:     derefOrEmpty(qSAMLConnection.IdpRedirectUrl),
		IdpX509Certificate: certPEM,
		IdpEntityId:        derefOrEmpty(qSAMLConnection.IdpEntityID),
	}, nil
}
