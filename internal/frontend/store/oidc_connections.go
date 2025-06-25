package store

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListOIDCConnections(ctx context.Context, req *frontendv1.ListOIDCConnectionsRequest) (*frontendv1.ListOIDCConnectionsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qOIDCConnections, err := q.ListOIDCConnections(ctx, queries.ListOIDCConnectionsParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list oidc connections: %w", err)
	}

	var oidcConnections []*frontendv1.OIDCConnection
	for _, qOIDCConn := range qOIDCConnections {
		oidcConnections = append(oidcConnections, parseOIDCConnection(qProject, qOIDCConn))
	}

	var nextPageToken string
	if len(oidcConnections) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qOIDCConnections[limit].ID)
		oidcConnections = oidcConnections[:limit]
	}

	return &frontendv1.ListOIDCConnectionsResponse{
		OidcConnections: oidcConnections,
		NextPageToken:   nextPageToken,
	}, nil
}

func (s *Store) GetOIDCConnection(ctx context.Context, req *frontendv1.GetOIDCConnectionRequest) (*frontendv1.GetOIDCConnectionResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	oidcConnectionID, err := idformat.OIDCConnection.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid oidc connection id", fmt.Errorf("parse oidc connection id: %w", err))
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOIDCConnection, err := q.GetOIDCConnection(ctx, queries.GetOIDCConnectionParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             oidcConnectionID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("oidc connection not found", fmt.Errorf("get oidc connection: %w", err))
		}

		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	return &frontendv1.GetOIDCConnectionResponse{OidcConnection: parseOIDCConnection(qProject, qOIDCConnection)}, nil
}

func (s *Store) CreateOIDCConnection(ctx context.Context, req *frontendv1.CreateOIDCConnectionRequest) (*frontendv1.CreateOIDCConnectionResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", err)
		}

		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	if !qOrg.LogInWithOidc {
		return nil, apierror.NewFailedPreconditionError("oidc is not enabled for the organization", fmt.Errorf("oidc is not enabled for the organization"))
	}

	if req.OidcConnection.ConfigurationUrl != "" {
		u, err := url.Parse(req.OidcConnection.ConfigurationUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid configuration url: %w", err)
		}

		if !u.IsAbs() {
			return nil, apierror.NewFailedPreconditionError("invalid idp configuration url", fmt.Errorf("invalid idp configuration url"))
		}

		config, err := s.oidc.GetConfiguration(ctx, req.OidcConnection.ConfigurationUrl)
		if err != nil {
			return nil, fmt.Errorf("get OIDC configuration: %w", err)
		}

		if err := config.Validate(); err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid oidc configuration", fmt.Errorf("validate oidc configuration: %w", err))
		}
	}

	var clientSecretCiphertext []byte
	if req.OidcConnection.ClientSecret != "" {
		encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               &s.oidcClientSecretsKMSKeyID,
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
			Plaintext:           []byte(req.OidcConnection.ClientSecret),
		})
		if err != nil {
			return nil, fmt.Errorf("encrypt oidc client secret: %w", err)
		}
		clientSecretCiphertext = encryptRes.CiphertextBlob
	}

	qOIDCConnection, err := q.CreateOIDCConnection(ctx, queries.CreateOIDCConnectionParams{
		ID:                     uuid.New(),
		OrganizationID:         authn.OrganizationID(ctx),
		IsPrimary:              derefOrEmpty(req.OidcConnection.Primary),
		ConfigurationUrl:       req.OidcConnection.ConfigurationUrl,
		Issuer:                 req.OidcConnection.Issuer,
		ClientID:               req.OidcConnection.ClientId,
		ClientSecretCiphertext: clientSecretCiphertext,
	})
	if err != nil {
		return nil, fmt.Errorf("create oidc connection: %w", err)
	}

	if req.OidcConnection.GetPrimary() {
		if err := q.UpdatePrimaryOIDCConnection(ctx, queries.UpdatePrimaryOIDCConnectionParams{
			OrganizationID: authn.OrganizationID(ctx),
			ID:             qOIDCConnection.ID,
		}); err != nil {
			return nil, fmt.Errorf("update primary oidc connection: %w", err)
		}
	}

	auditOIDCConnection, err := s.auditlogStore.GetOIDCConnection(ctx, tx, qOIDCConnection.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit oidc connection: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.oidc_connections.create",
		EventDetails: &auditlogv1.CreateOIDCConnection{
			OidcConnection: auditOIDCConnection,
		},
		ResourceType: queries.AuditLogEventResourceTypeOidcConnection,
		ResourceID:   &qOIDCConnection.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.CreateOIDCConnectionResponse{OidcConnection: parseOIDCConnection(qProject, qOIDCConnection)}, nil
}

func (s *Store) UpdateOIDCConnection(ctx context.Context, req *frontendv1.UpdateOIDCConnectionRequest) (*frontendv1.UpdateOIDCConnectionResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	oidcConnectionID, err := idformat.OIDCConnection.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid oidc connection id", fmt.Errorf("parse oidc connection id: %w", err))
	}

	// authz
	qOIDCConnection, err := q.GetOIDCConnection(ctx, queries.GetOIDCConnectionParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             oidcConnectionID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("oidc connection not found", fmt.Errorf("get oidc connection: %w", err))
		}

		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	auditPreviousOIDCConnection, err := s.auditlogStore.GetOIDCConnection(ctx, tx, qOIDCConnection.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit previous oidc connection: %w", err)
	}

	updates := queries.UpdateOIDCConnectionParams{
		ID:                     oidcConnectionID,
		IsPrimary:              qOIDCConnection.IsPrimary,
		ConfigurationUrl:       qOIDCConnection.ConfigurationUrl,
		Issuer:                 qOIDCConnection.Issuer,
		ClientID:               qOIDCConnection.ClientID,
		ClientSecretCiphertext: qOIDCConnection.ClientSecretCiphertext,
	}

	if req.OidcConnection.ConfigurationUrl != "" {
		u, err := url.Parse(req.OidcConnection.ConfigurationUrl)
		if err != nil {
			return nil, apierror.NewFailedPreconditionError("invalid configuration url", fmt.Errorf("invalid configuration url: %w", err))
		}

		if !u.IsAbs() {
			return nil, apierror.NewFailedPreconditionError("invalid configuration url", fmt.Errorf("invalid configuration url"))
		}

		config, err := s.oidc.GetConfiguration(ctx, req.OidcConnection.ConfigurationUrl)
		if err != nil {
			return nil, fmt.Errorf("get OIDC configuration: %w", err)
		}

		if err := config.Validate(); err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid oidc configuration", fmt.Errorf("validate oidc configuration: %w", err))
		}

		updates.ConfigurationUrl = req.OidcConnection.ConfigurationUrl
	}

	if req.OidcConnection.Issuer != "" {
		if _, err := url.Parse(req.OidcConnection.Issuer); err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid oidc issuer", fmt.Errorf("invalid oidc issuer: %w", err))
		}
		updates.Issuer = req.OidcConnection.Issuer
	}

	if req.OidcConnection.ClientId != "" {
		updates.ClientID = req.OidcConnection.ClientId
	}

	if req.OidcConnection.ClientSecret != "" {
		encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               &s.oidcClientSecretsKMSKeyID,
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
			Plaintext:           []byte(req.OidcConnection.ClientSecret),
		})
		if err != nil {
			return nil, fmt.Errorf("encrypt oidc client secret: %w", err)
		}
		updates.ClientSecretCiphertext = encryptRes.CiphertextBlob
	}

	if req.OidcConnection.Primary != nil {
		updates.IsPrimary = *req.OidcConnection.Primary
	}

	qUpdatedOIDCConnection, err := q.UpdateOIDCConnection(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update oidc connection: %w", err)
	}

	auditOIDCConnection, err := s.auditlogStore.GetOIDCConnection(ctx, tx, qUpdatedOIDCConnection.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit oidc connection: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.oidc_connections.update",
		EventDetails: &auditlogv1.UpdateOIDCConnection{
			OidcConnection:         auditOIDCConnection,
			PreviousOidcConnection: auditPreviousOIDCConnection,
		},
		ResourceType: queries.AuditLogEventResourceTypeOidcConnection,
		ResourceID:   &qOIDCConnection.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.UpdateOIDCConnectionResponse{OidcConnection: parseOIDCConnection(qProject, qUpdatedOIDCConnection)}, nil
}

func (s *Store) DeleteOIDCConnection(ctx context.Context, req *frontendv1.DeleteOIDCConnectionRequest) (*frontendv1.DeleteOIDCConnectionResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	oidcConnectionID, err := idformat.OIDCConnection.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid oidc connection id", fmt.Errorf("parse oidc connection id: %w", err))
	}

	// authz
	qOIDCConnection, err := q.GetOIDCConnection(ctx, queries.GetOIDCConnectionParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             oidcConnectionID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("oidc connection not found", fmt.Errorf("get oidc connection: %w", err))
		}

		return nil, fmt.Errorf("get oidc connection: %w", err)
	}

	auditOIDCConnection, err := s.auditlogStore.GetOIDCConnection(ctx, tx, qOIDCConnection.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit oidc connection: %w", err)
	}

	if err := q.DeleteOIDCConnection(ctx, oidcConnectionID); err != nil {
		return nil, fmt.Errorf("delete oidc connection: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.oidc_connections.delete",
		EventDetails: &auditlogv1.DeleteOIDCConnection{
			OidcConnection: auditOIDCConnection,
		},
		ResourceType: queries.AuditLogEventResourceTypeOidcConnection,
		ResourceID:   &qOIDCConnection.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.DeleteOIDCConnectionResponse{}, nil
}

func parseOIDCConnection(qProject queries.Project, qOIDCConnection queries.OidcConnection) *frontendv1.OIDCConnection {
	redirectURL := fmt.Sprintf("https://%s/api/oidc/v1/%s/callback", qProject.VaultDomain, idformat.OIDCConnection.Format(qOIDCConnection.ID))

	return &frontendv1.OIDCConnection{
		Id:               idformat.OIDCConnection.Format(qOIDCConnection.ID),
		CreateTime:       timestamppb.New(*qOIDCConnection.CreateTime),
		UpdateTime:       timestamppb.New(*qOIDCConnection.UpdateTime),
		Primary:          &qOIDCConnection.IsPrimary,
		Issuer:           qOIDCConnection.Issuer,
		ConfigurationUrl: qOIDCConnection.ConfigurationUrl,
		ClientId:         qOIDCConnection.ClientID,
		ClientSecret:     "",
		RedirectUri:      redirectURL,
	}
}
