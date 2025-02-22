package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ListPublishableKeys(ctx context.Context, req *connect.Request[backendv1.ListPublishableKeysRequest]) (*connect.Response[backendv1.ListPublishableKeysResponse], error) {
	res, err := s.Store.ListPublishableKeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetPublishableKey(ctx context.Context, req *connect.Request[backendv1.GetPublishableKeyRequest]) (*connect.Response[backendv1.GetPublishableKeyResponse], error) {
	res, err := s.Store.GetPublishableKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreatePublishableKey(ctx context.Context, req *connect.Request[backendv1.CreatePublishableKeyRequest]) (*connect.Response[backendv1.CreatePublishableKeyResponse], error) {
	res, err := s.Store.CreatePublishableKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdatePublishableKey(ctx context.Context, req *connect.Request[backendv1.UpdatePublishableKeyRequest]) (*connect.Response[backendv1.UpdatePublishableKeyResponse], error) {
	res, err := s.Store.UpdatePublishableKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeletePublishableKey(ctx context.Context, req *connect.Request[backendv1.DeletePublishableKeyRequest]) (*connect.Response[backendv1.DeletePublishableKeyResponse], error) {
	res, err := s.Store.DeletePublishableKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
