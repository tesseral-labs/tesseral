package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
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

func (s *Service) CreateUser(ctx context.Context, req *connect.Request[backendv1.CreateUserRequest]) (*connect.Response[backendv1.CreateUserResponse], error) {
	res, err := s.Store.CreateUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) UpdateUser(ctx context.Context, req *connect.Request[backendv1.UpdateUserRequest]) (*connect.Response[backendv1.UpdateUserResponse], error) {
	res, err := s.Store.UpdateUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteUser(ctx context.Context, req *connect.Request[backendv1.DeleteUserRequest]) (*connect.Response[backendv1.DeleteUserResponse], error) {
	res, err := s.Store.DeleteUser(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
