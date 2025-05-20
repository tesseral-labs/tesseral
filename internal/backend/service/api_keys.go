package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) CreateAPIKey(ctx context.Context, req *connect.Request[backendv1.CreateAPIKeyRequest]) (*connect.Response[backendv1.CreateAPIKeyResponse], error) {
	res, err := s.Store.CreateAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteAPIKey(ctx context.Context, req *connect.Request[backendv1.DeleteAPIKeyRequest]) (*connect.Response[backendv1.DeleteAPIKeyResponse], error) {
	res, err := s.Store.DeleteAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetAPIKey(ctx context.Context, req *connect.Request[backendv1.GetAPIKeyRequest]) (*connect.Response[backendv1.GetAPIKeyResponse], error) {
	res, err := s.Store.GetAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListAPIKeys(ctx context.Context, req *connect.Request[backendv1.ListAPIKeysRequest]) (*connect.Response[backendv1.ListAPIKeysResponse], error) {
	res, err := s.Store.ListAPIKeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeAPIKey(ctx context.Context, req *connect.Request[backendv1.RevokeAPIKeyRequest]) (*connect.Response[backendv1.RevokeAPIKeyResponse], error) {
	res, err := s.Store.RevokeAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateAPIKey(ctx context.Context, req *connect.Request[backendv1.UpdateAPIKeyRequest]) (*connect.Response[backendv1.UpdateAPIKeyResponse], error) {
	res, err := s.Store.UpdateAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) AuthenticateAPIKey(ctx context.Context, req *connect.Request[backendv1.AuthenticateAPIKeyRequest]) (*connect.Response[backendv1.AuthenticateAPIKeyResponse], error) {
	res, err := s.Store.AuthenticateAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
