package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) ListUserInvites(ctx context.Context, req *connect.Request[backendv1.ListUserInvitesRequest]) (*connect.Response[backendv1.ListUserInvitesResponse], error) {
	res, err := s.Store.ListUserInvites(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetUserInvite(ctx context.Context, req *connect.Request[backendv1.GetUserInviteRequest]) (*connect.Response[backendv1.GetUserInviteResponse], error) {
	data, err := s.Store.GetUserInvite(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(data), nil
}

func (s *Service) CreateUserInvite(ctx context.Context, req *connect.Request[backendv1.CreateUserInviteRequest]) (*connect.Response[backendv1.CreateUserInviteResponse], error) {
	data, err := s.Store.CreateUserInvite(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(data), nil
}

func (s *Service) DeleteUserInvite(ctx context.Context, req *connect.Request[backendv1.DeleteUserInviteRequest]) (*connect.Response[backendv1.DeleteUserInviteResponse], error) {
	res, err := s.Store.DeleteUserInvite(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
