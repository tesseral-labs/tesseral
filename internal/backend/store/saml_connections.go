package store

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListSAMLConnections(ctx context.Context, req *backendv1.ListSAMLConnectionsRequest) (*backendv1.ListSAMLConnectionsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.OrganizationId)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	// authz
	if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        orgID,
	}); err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qSAMLConnections, err := q.ListSAMLConnections(ctx, queries.ListSAMLConnectionsParams{
		OrganizationID: orgID,
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list saml connections: %w", err)
	}

	var samlConnections []*backendv1.SAMLConnection
	for _, qSAMLConn := range qSAMLConnections {
		samlConnections = append(samlConnections, parseSAMLConnection(qSAMLConn))
	}

	var nextPageToken string
	if len(samlConnections) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qSAMLConnections[limit].ID)
		samlConnections = samlConnections[:limit]
	}

	return &backendv1.ListSAMLConnectionsResponse{
		SamlConnections: samlConnections,
		NextPageToken:   nextPageToken,
	}, nil
}

func (s *Store) GetSAMLConnection(ctx context.Context, req *backendv1.GetSAMLConnectionRequest) (*backendv1.GetSAMLConnectionResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	samlConnectionID, err := idformat.SAMLConnection.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse saml connection id: %w", err)
	}

	qSAMLConnection, err := q.GetSAMLConnection(ctx, queries.GetSAMLConnectionParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        samlConnectionID,
	})
	if err != nil {
		return nil, fmt.Errorf("get saml connection: %w", err)
	}

	return &backendv1.GetSAMLConnectionResponse{SamlConnection: parseSAMLConnection(qSAMLConnection)}, nil
}

func (s *Store) CreateSAMLConnection(ctx context.Context, req *backendv1.CreateSAMLConnectionRequest) (*backendv1.CreateSAMLConnectionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.SamlConnection.OrganizationId)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	// authz
	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	if !qOrg.SamlEnabled {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("organization does not have SAML enabled"))
	}

	if req.SamlConnection.IdpRedirectUrl != "" {
		u, err := url.Parse(req.SamlConnection.IdpRedirectUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid idp redirect url: %w", err)
		}

		if !u.IsAbs() {
			return nil, fmt.Errorf("idp redirect url must be absolute")
		}
	}

	var idpCertificate []byte
	if req.SamlConnection.IdpX509Certificate != "" {
		block, _ := pem.Decode([]byte(req.SamlConnection.IdpX509Certificate))
		if block == nil || block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("invalid certificate format")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %w", err)
		}

		idpCertificate = cert.Raw
	}

	qSAMLConnection, err := q.CreateSAMLConnection(ctx, queries.CreateSAMLConnectionParams{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		IsPrimary:          derefOrEmpty(req.SamlConnection.Primary),
		IdpRedirectUrl:     &req.SamlConnection.IdpRedirectUrl,
		IdpX509Certificate: idpCertificate,
		IdpEntityID:        &req.SamlConnection.IdpEntityId,
	})
	if err != nil {
		return nil, fmt.Errorf("create saml connection: %w", err)
	}

	if req.SamlConnection.GetPrimary() {
		if err := q.UpdatePrimarySAMLConnection(ctx, queries.UpdatePrimarySAMLConnectionParams{
			OrganizationID: orgID,
			ID:             qSAMLConnection.ID,
		}); err != nil {
			return nil, fmt.Errorf("update primary saml connection: %w", err)
		}
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.CreateSAMLConnectionResponse{SamlConnection: parseSAMLConnection(qSAMLConnection)}, nil
}

func (s *Store) UpdateSAMLConnection(ctx context.Context, req *backendv1.UpdateSAMLConnectionRequest) (*backendv1.UpdateSAMLConnectionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	samlConnectionID, err := idformat.SAMLConnection.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse saml connection id: %w", err)
	}

	// authz
	qSAMLConnection, err := q.GetSAMLConnection(ctx, queries.GetSAMLConnectionParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        samlConnectionID,
	})
	if err != nil {
		return nil, fmt.Errorf("get saml connection: %w", err)
	}

	updates := queries.UpdateSAMLConnectionParams{
		ID:                 samlConnectionID,
		IsPrimary:          qSAMLConnection.IsPrimary,
		IdpRedirectUrl:     qSAMLConnection.IdpRedirectUrl,
		IdpX509Certificate: qSAMLConnection.IdpX509Certificate,
		IdpEntityID:        qSAMLConnection.IdpEntityID,
	}

	if req.SamlConnection.IdpRedirectUrl != "" {
		u, err := url.Parse(req.SamlConnection.IdpRedirectUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid idp redirect url: %w", err)
		}

		if !u.IsAbs() {
			return nil, fmt.Errorf("idp redirect url must be absolute")
		}

		updates.IdpRedirectUrl = &req.SamlConnection.IdpRedirectUrl
	}

	if req.SamlConnection.IdpX509Certificate != "" {
		block, _ := pem.Decode([]byte(req.SamlConnection.IdpX509Certificate))
		if block == nil || block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("invalid certificate format")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %w", err)
		}

		updates.IdpX509Certificate = cert.Raw
	}

	if req.SamlConnection.IdpEntityId != "" {
		updates.IdpEntityID = &req.SamlConnection.IdpEntityId
	}

	if req.SamlConnection.Primary != nil {
		updates.IsPrimary = *req.SamlConnection.Primary
	}

	qUpdatedSAMLConnection, err := q.UpdateSAMLConnection(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update saml connection: %w", err)
	}

	if req.SamlConnection.GetPrimary() {
		if err := q.UpdatePrimarySAMLConnection(ctx, queries.UpdatePrimarySAMLConnectionParams{
			OrganizationID: qSAMLConnection.OrganizationID,
			ID:             samlConnectionID,
		}); err != nil {
			return nil, fmt.Errorf("update primary saml connection: %w", err)
		}
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateSAMLConnectionResponse{SamlConnection: parseSAMLConnection(qUpdatedSAMLConnection)}, nil
}

func (s *Store) DeleteSAMLConnection(ctx context.Context, req *backendv1.DeleteSAMLConnectionRequest) (*backendv1.DeleteSAMLConnectionResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	samlConnectionID, err := idformat.SAMLConnection.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse saml connection id: %w", err)
	}

	// authz
	if _, err := q.GetSAMLConnection(ctx, queries.GetSAMLConnectionParams{
		ProjectID: projectid.ProjectID(ctx),
		ID:        samlConnectionID,
	}); err != nil {
		return nil, fmt.Errorf("get saml connection: %w", err)
	}

	if err := q.DeleteSAMLConnection(ctx, samlConnectionID); err != nil {
		return nil, fmt.Errorf("delete saml connection: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteSAMLConnectionResponse{}, nil
}

func parseSAMLConnection(qSAMLConnection queries.SamlConnection) *backendv1.SAMLConnection {
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

	spACSURL := fmt.Sprintf("http://localhost:3001/saml/v1/%s", idformat.SAMLConnection.Format(qSAMLConnection.ID))   // todo
	spEntityID := fmt.Sprintf("http://localhost:3001/saml/v1/%s", idformat.SAMLConnection.Format(qSAMLConnection.ID)) // todo

	return &backendv1.SAMLConnection{
		Id:                 idformat.SAMLConnection.Format(qSAMLConnection.ID),
		OrganizationId:     idformat.Organization.Format(qSAMLConnection.OrganizationID),
		CreateTime:         timestamppb.New(*qSAMLConnection.CreateTime),
		UpdateTime:         timestamppb.New(*qSAMLConnection.UpdateTime),
		Primary:            &qSAMLConnection.IsPrimary,
		SpAcsUrl:           spACSURL,
		SpEntityId:         spEntityID,
		IdpRedirectUrl:     derefOrEmpty(qSAMLConnection.IdpRedirectUrl),
		IdpX509Certificate: certPEM,
		IdpEntityId:        derefOrEmpty(qSAMLConnection.IdpEntityID),
	}
}
