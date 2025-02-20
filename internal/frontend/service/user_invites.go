package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) ListUserInvites(ctx context.Context, req *connect.Request[frontendv1.ListUserInvitesRequest]) (*connect.Response[frontendv1.ListUserInvitesResponse], error) {
	res, err := s.Store.ListUserInvites(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetUserInvite(ctx context.Context, req *connect.Request[frontendv1.GetUserInviteRequest]) (*connect.Response[frontendv1.GetUserInviteResponse], error) {
	data, err := s.Store.GetUserInvite(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(data), nil
}

func (s *Service) CreateUserInvite(ctx context.Context, req *connect.Request[frontendv1.CreateUserInviteRequest]) (*connect.Response[frontendv1.CreateUserInviteResponse], error) {
	data, err := s.Store.CreateUserInvite(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(data), nil
}

func (s *Service) DeleteUserInvite(ctx context.Context, req *connect.Request[frontendv1.DeleteUserInviteRequest]) (*connect.Response[frontendv1.DeleteUserInviteResponse], error) {
	res, err := s.Store.DeleteUserInvite(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
