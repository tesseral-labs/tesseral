package store

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/frontend/authn"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/webauthn"
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
		RpId:            *qProject.AuthDomain,
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
		RPID:              *qProject.AuthDomain,
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
	})
	if err != nil {
		return nil, fmt.Errorf("create passkey: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.RegisterPasskeyResponse{
		Passkey: parsePasskey(qPasskey),
	}, nil
}

func parsePasskey(qPasskey queries.Passkey) *frontendv1.Passkey {
	return &frontendv1.Passkey{
		Id:           idformat.Passkey.Format(qPasskey.ID),
		UserId:       idformat.User.Format(qPasskey.UserID),
		CreateTime:   timestamppb.New(*qPasskey.CreateTime),
		UpdateTime:   timestamppb.New(*qPasskey.UpdateTime),
		CredentialId: qPasskey.CredentialID,
		PublicKeyPkix: string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: qPasskey.PublicKey,
		})),
		Aaguid: qPasskey.Aaguid,
	}
}
