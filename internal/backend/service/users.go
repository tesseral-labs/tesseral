package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) ListUsers(ctx context.Context, req *connect.Request[backendv1.ListUsersRequest]) (*connect.Response[backendv1.ListUsersResponse], error) {
	res, err := s.Store.ListUsers(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetUser(ctx context.Context, req *connect.Request[backendv1.GetUserRequest]) (*connect.Response[backendv1.GetUserResponse], error) {
	res, err := s.Store.GetUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
