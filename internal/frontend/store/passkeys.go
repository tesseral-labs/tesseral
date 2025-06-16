package store

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/webauthn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListMyPasskeys(ctx context.Context, req *frontendv1.ListMyPasskeysRequest) (*frontendv1.ListMyPasskeysResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qPasskeys, err := q.ListPasskeys(ctx, queries.ListPasskeysParams{
		UserID: authn.UserID(ctx),
		ID:     startID,
		Limit:  int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list passkeys: %w", err)
	}

	var passkeys []*frontendv1.Passkey
	for _, qPasskey := range qPasskeys {
		passkeys = append(passkeys, parsePasskey(qPasskey))
	}

	var nextPageToken string
	if len(passkeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qPasskeys[limit].ID)
		passkeys = passkeys[:limit]
	}

	return &frontendv1.ListMyPasskeysResponse{
		Passkeys:      passkeys,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) DeleteMyPasskey(ctx context.Context, req *frontendv1.DeleteMyPasskeyRequest) (*frontendv1.DeleteMyPasskeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	passkeyID, err := idformat.Passkey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse passkey id: %w", err)
	}

	qPasskey, err := q.GetUserPasskey(ctx, queries.GetUserPasskeyParams{
		UserID: authn.UserID(ctx),
		ID:     passkeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("passkey not found", fmt.Errorf("get user passkey: %w", err))
		}

		return nil, fmt.Errorf("get user passkey: %w", err)
	}

	if err := q.DeletePasskey(ctx, passkeyID); err != nil {
		return nil, fmt.Errorf("delete passkey: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.passkeys.delete",
		EventDetails: &frontendv1.PasskeyDeleted{
			Passkey: parsePasskey(qPasskey),
		},
		ResourceType: queries.AuditLogEventResourceTypePasskey,
		ResourceID:   &qPasskey.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.DeleteMyPasskeyResponse{}, nil
}

func (s *Store) GetPasskeyOptions(ctx context.Context, req *frontendv1.GetPasskeyOptionsRequest) (*frontendv1.GetPasskeyOptionsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qUser, err := q.GetUserByID(ctx, authn.UserID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &frontendv1.GetPasskeyOptionsResponse{
		RpId:            qProject.VaultDomain,
		RpName:          qProject.DisplayName,
		UserId:          fmt.Sprintf("%s|%s", qUser.Email, idformat.Organization.Format(qUser.OrganizationID)),
		UserDisplayName: qUser.Email,
	}, nil
}

func (s *Store) RegisterPasskey(ctx context.Context, req *frontendv1.RegisterPasskeyRequest) (*frontendv1.RegisterPasskeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	cred, err := webauthn.Parse(&webauthn.ParseRequest{
		RPID:              qProject.VaultDomain,
		AttestationObject: req.AttestationObject,
	})
	if err != nil {
		return nil, fmt.Errorf("parse webauthn credential: %w", err)
	}

	publicKey, err := x509.MarshalPKIXPublicKey(cred.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	qPasskey, err := q.CreatePasskey(ctx, queries.CreatePasskeyParams{
		ID:           uuid.New(),
		UserID:       authn.UserID(ctx),
		CredentialID: cred.ID,
		PublicKey:    publicKey,
		Aaguid:       cred.AAGUID,
		RpID:         req.RpId,
	})
	if err != nil {
		return nil, fmt.Errorf("create passkey: %w", err)
	}

	passkey := parsePasskey(qPasskey)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.passkeys.create",
		EventDetails: &frontendv1.PasskeyCreated{
			Passkey: passkey,
		},
		ResourceType: queries.AuditLogEventResourceTypePasskey,
		ResourceID:   &qPasskey.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.RegisterPasskeyResponse{
		Passkey: passkey,
	}, nil
}

func parsePasskey(qPasskey queries.Passkey) *frontendv1.Passkey {
	return &frontendv1.Passkey{
		Id:           idformat.Passkey.Format(qPasskey.ID),
		UserId:       idformat.User.Format(qPasskey.UserID),
		CreateTime:   timestamppb.New(*qPasskey.CreateTime),
		UpdateTime:   timestamppb.New(*qPasskey.UpdateTime),
		Disabled:     qPasskey.Disabled,
		CredentialId: qPasskey.CredentialID,
		PublicKeyPkix: string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: qPasskey.PublicKey,
		})),
		Aaguid: qPasskey.Aaguid,
		RpId:   qPasskey.RpID,
	}
}
