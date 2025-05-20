package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) CreateAPIKey(ctx context.Context, req *connect.Request[frontendv1.CreateAPIKeyRequest]) (*connect.Response[frontendv1.CreateAPIKeyResponse], error) {
	res, err := s.Store.CreateAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteAPIKey(ctx context.Context, req *connect.Request[frontendv1.DeleteAPIKeyRequest]) (*connect.Response[frontendv1.DeleteAPIKeyResponse], error) {
	res, err := s.Store.DeleteAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetAPIKey(ctx context.Context, req *connect.Request[frontendv1.GetAPIKeyRequest]) (*connect.Response[frontendv1.GetAPIKeyResponse], error) {
	res, err := s.Store.GetAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListAPIKeys(ctx context.Context, req *connect.Request[frontendv1.ListAPIKeysRequest]) (*connect.Response[frontendv1.ListAPIKeysResponse], error) {
	res, err := s.Store.ListAPIKeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeAPIKey(ctx context.Context, req *connect.Request[frontendv1.RevokeAPIKeyRequest]) (*connect.Response[frontendv1.RevokeAPIKeyResponse], error) {
	res, err := s.Store.RevokeAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateAPIKey(ctx context.Context, req *connect.Request[frontendv1.UpdateAPIKeyRequest]) (*connect.Response[frontendv1.UpdateAPIKeyResponse], error) {
	res, err := s.Store.UpdateAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
