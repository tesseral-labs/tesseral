package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) ListProjectAPIKeys(ctx context.Context, req *connect.Request[backendv1.ListProjectAPIKeysRequest]) (*connect.Response[backendv1.ListProjectAPIKeysResponse], error) {
	res, err := s.Store.ListProjectAPIKeys(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetProjectAPIKey(ctx context.Context, req *connect.Request[backendv1.GetProjectAPIKeyRequest]) (*connect.Response[backendv1.GetProjectAPIKeyResponse], error) {
	res, err := s.Store.GetProjectAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) CreateProjectAPIKey(ctx context.Context, req *connect.Request[backendv1.CreateProjectAPIKeyRequest]) (*connect.Response[backendv1.CreateProjectAPIKeyResponse], error) {
	res, err := s.Store.CreateProjectAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateProjectAPIKey(ctx context.Context, req *connect.Request[backendv1.UpdateProjectAPIKeyRequest]) (*connect.Response[backendv1.UpdateProjectAPIKeyResponse], error) {
	res, err := s.Store.UpdateProjectAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteProjectAPIKey(ctx context.Context, req *connect.Request[backendv1.DeleteProjectAPIKeyRequest]) (*connect.Response[backendv1.DeleteProjectAPIKeyResponse], error) {
	res, err := s.Store.DeleteProjectAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) RevokeProjectAPIKey(ctx context.Context, req *connect.Request[backendv1.RevokeProjectAPIKeyRequest]) (*connect.Response[backendv1.RevokeProjectAPIKeyResponse], error) {
	res, err := s.Store.RevokeProjectAPIKey(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
