package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) ListUsers(ctx context.Context, req *connect.Request[frontendv1.ListUsersRequest]) (*connect.Response[frontendv1.ListUsersResponse], error) {
	res, err := s.Store.ListUsers(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetUser(ctx context.Context, req *connect.Request[frontendv1.GetUserRequest]) (*connect.Response[frontendv1.GetUserResponse], error) {
	res, err := s.Store.GetUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateUser(ctx context.Context, req *connect.Request[frontendv1.UpdateUserRequest]) (*connect.Response[frontendv1.UpdateUserResponse], error) {
	res, err := s.Store.UpdateUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
