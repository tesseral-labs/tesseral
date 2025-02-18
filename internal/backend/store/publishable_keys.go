package store

import (
	"context"
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

func (s *Store) ListPublishableKeys(ctx context.Context, req *backendv1.ListPublishableKeysRequest) (*backendv1.ListPublishableKeysResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	qPublishableKeys, err := q.ListPublishableKeys(ctx, queries.ListPublishableKeysParams{
		ProjectID: authn.ProjectID(ctx),
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list publishable keys: %w", err)
	}

	var publishableKeys []*backendv1.PublishableKey
	for _, qPublishableKey := range qPublishableKeys {
		publishableKeys = append(publishableKeys, parsePublishableKey(qPublishableKey))
	}

	var nextPageToken string
	if len(publishableKeys) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(publishableKeys[limit].Id)
		publishableKeys = publishableKeys[:limit]
	}

	return &backendv1.ListPublishableKeysResponse{
		PublishableKeys: publishableKeys,
		NextPageToken:   nextPageToken,
	}, nil
}

func (s *Store) GetPublishableKey(ctx context.Context, req *backendv1.GetPublishableKeyRequest) (*backendv1.GetPublishableKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	id, err := idformat.PublishableKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid publishable key id", fmt.Errorf("parse publishable key id: %w", err))
	}

	qPublishableKey, err := s.q.GetPublishableKey(ctx, queries.GetPublishableKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("publishable key not found", fmt.Errorf("get publishable key: %w", err))
		}

		return nil, fmt.Errorf("get publishable key: %w", err)
	}

	return &backendv1.GetPublishableKeyResponse{PublishableKey: parsePublishableKey(qPublishableKey)}, nil
}

func (s *Store) CreatePublishableKey(ctx context.Context, req *backendv1.CreatePublishableKeyRequest) (*backendv1.CreatePublishableKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qPublishableKey, err := q.CreatePublishableKey(ctx, queries.CreatePublishableKeyParams{
		ID:          uuid.New(),
		ProjectID:   authn.ProjectID(ctx),
		DisplayName: req.PublishableKey.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("create publishable key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.CreatePublishableKeyResponse{PublishableKey: parsePublishableKey(qPublishableKey)}, nil
}

func (s *Store) UpdatePublishableKey(ctx context.Context, req *backendv1.UpdatePublishableKeyRequest) (*backendv1.UpdatePublishableKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	publishableKeyID, err := idformat.PublishableKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid publishable key id", fmt.Errorf("parse publishable key id: %w", err))
	}

	qPublishableKey, err := q.GetPublishableKey(ctx, queries.GetPublishableKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        publishableKeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("publishable key not found", fmt.Errorf("get publishable key: %w", err))
		}

		return nil, fmt.Errorf("get publishable key: %w", err)
	}

	updates := queries.UpdatePublishableKeyParams{
		ID:          publishableKeyID,
		DisplayName: qPublishableKey.DisplayName,
	}

	if req.PublishableKey.DisplayName != "" {
		updates.DisplayName = req.PublishableKey.DisplayName
	}

	qUpdatedPublishableKey, err := q.UpdatePublishableKey(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update publishable key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdatePublishableKeyResponse{PublishableKey: parsePublishableKey(qUpdatedPublishableKey)}, nil
}

func (s *Store) DeletePublishableKey(ctx context.Context, req *backendv1.DeletePublishableKeyRequest) (*backendv1.DeletePublishableKeyResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	publishableKeyID, err := idformat.PublishableKey.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid publishable key id", fmt.Errorf("parse publishable key id: %w", err))
	}

	if _, err := q.GetPublishableKey(ctx, queries.GetPublishableKeyParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        publishableKeyID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("publishable key not found", fmt.Errorf("get publishable key: %w", err))
		}

		return nil, fmt.Errorf("get publishable key: %w", err)
	}

	if err := q.DeletePublishableKey(ctx, publishableKeyID); err != nil {
		return nil, fmt.Errorf("delete publishable key: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeletePublishableKeyResponse{}, nil
}

func parsePublishableKey(qPublishableKey queries.PublishableKey) *backendv1.PublishableKey {
	return &backendv1.PublishableKey{
		Id:          idformat.PublishableKey.Format(qPublishableKey.ID),
		DisplayName: qPublishableKey.DisplayName,
		CreateTime:  timestamppb.New(*qPublishableKey.CreateTime),
		UpdateTime:  timestamppb.New(*qPublishableKey.UpdateTime),
	}
}
