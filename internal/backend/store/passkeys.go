package store

import (
	"context"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListPasskeys(ctx context.Context, req *backendv1.ListPasskeysRequest) (*backendv1.ListPasskeysResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	if _, err := q.GetUser(ctx, queries.GetUserParams{
		ID:        userID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qPasskeys, err := q.ListPasskeys(ctx, queries.ListPasskeysParams{
		UserID: userID,
		ID:     startID,
		Limit:  int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list passkeys: %w", err)
	}

	var passkeys []*backendv1.Passkey
	for _, qPasskey := range qPasskeys {
		passkeys = append(passkeys, parsePasskey(qPasskey))
	}

	var nextPageToken string
	if len(passkeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qPasskeys[limit].ID)
		passkeys = passkeys[:limit]
	}

	return &backendv1.ListPasskeysResponse{
		Passkeys:      passkeys,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetPasskey(ctx context.Context, req *backendv1.GetPasskeyRequest) (*backendv1.GetPasskeyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	passkeyID, err := idformat.Passkey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse passkey id: %w", err)
	}

	qPasskey, err := q.GetPasskey(ctx, queries.GetPasskeyParams{
		ID:        passkeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("passkey not found", fmt.Errorf("get passkey: %w", err))
		}

		return nil, fmt.Errorf("get passkey: %w", err)
	}

	return &backendv1.GetPasskeyResponse{Passkey: parsePasskey(qPasskey)}, nil
}

func (s *Store) UpdatePasskey(ctx context.Context, req *backendv1.UpdatePasskeyRequest) (*backendv1.UpdatePasskeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	passkeyID, err := idformat.Passkey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse passkey id: %w", err)
	}

	qPasskey, err := q.GetPasskey(ctx, queries.GetPasskeyParams{
		ID:        passkeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("passkey not found", fmt.Errorf("get passkey: %w", err))
		}

		return nil, fmt.Errorf("get passkey: %w", err)
	}

	var updates queries.UpdatePasskeyParams
	updates.ID = qPasskey.ID

	if req.Passkey.Disabled != nil {
		updates.Disabled = *req.Passkey.Disabled
	}

	qUpdatedPasskey, err := q.UpdatePasskey(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update passkey: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdatePasskeyResponse{Passkey: parsePasskey(qUpdatedPasskey)}, nil
}

func (s *Store) DeletePasskey(ctx context.Context, req *backendv1.DeletePasskeyRequest) (*backendv1.DeletePasskeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	passkeyID, err := idformat.Passkey.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse passkey id: %w", err)
	}

	if _, err := q.GetPasskey(ctx, queries.GetPasskeyParams{
		ID:        passkeyID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("passkey not found", fmt.Errorf("get passkey: %w", err))
		}

		return nil, fmt.Errorf("get passkey: %w", err)
	}

	if err := q.DeletePasskey(ctx, passkeyID); err != nil {
		return nil, fmt.Errorf("delete passkey: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeletePasskeyResponse{}, nil
}

func parsePasskey(qPasskey queries.Passkey) *backendv1.Passkey {
	return &backendv1.Passkey{
		Id:           idformat.Passkey.Format(qPasskey.ID),
		UserId:       idformat.User.Format(qPasskey.UserID),
		CreateTime:   timestamppb.New(*qPasskey.CreateTime),
		UpdateTime:   timestamppb.New(*qPasskey.UpdateTime),
		Disabled:     &qPasskey.Disabled,
		CredentialId: qPasskey.CredentialID,
		PublicKeyPkix: string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: qPasskey.PublicKey,
		})),
		Aaguid: qPasskey.Aaguid,
		RpId:   qPasskey.RpID,
	}
}
