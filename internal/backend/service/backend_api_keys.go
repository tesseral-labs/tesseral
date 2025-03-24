package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ListBackendAPIKeys(ctx context.Context, req *connect.Request[backendv1.ListBackendAPIKeysRequest]) (*connect.Response[backendv1.ListBackendAPIKeysResponse], error) {
	res, err := s.Store.ListBackendAPIKeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetBackendAPIKey(ctx context.Context, req *connect.Request[backendv1.GetBackendAPIKeyRequest]) (*connect.Response[backendv1.GetBackendAPIKeyResponse], error) {
	res, err := s.Store.GetBackendAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateBackendAPIKey(ctx context.Context, req *connect.Request[backendv1.CreateBackendAPIKeyRequest]) (*connect.Response[backendv1.CreateBackendAPIKeyResponse], error) {
	res, err := s.Store.CreateBackendAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateBackendAPIKey(ctx context.Context, req *connect.Request[backendv1.UpdateBackendAPIKeyRequest]) (*connect.Response[backendv1.UpdateBackendAPIKeyResponse], error) {
	res, err := s.Store.UpdateBackendAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteBackendAPIKey(ctx context.Context, req *connect.Request[backendv1.DeleteBackendAPIKeyRequest]) (*connect.Response[backendv1.DeleteBackendAPIKeyResponse], error) {
	res, err := s.Store.DeleteBackendAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeBackendAPIKey(ctx context.Context, req *connect.Request[backendv1.RevokeBackendAPIKeyRequest]) (*connect.Response[backendv1.RevokeBackendAPIKeyResponse], error) {
	res, err := s.Store.RevokeBackendAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
