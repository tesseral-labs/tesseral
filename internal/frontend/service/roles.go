package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) ListRoles(ctx context.Context, req *connect.Request[frontendv1.ListRolesRequest]) (*connect.Response[frontendv1.ListRolesResponse], error) {
	res, err := s.Store.ListRoles(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetRole(ctx context.Context, req *connect.Request[frontendv1.GetRoleRequest]) (*connect.Response[frontendv1.GetRoleResponse], error) {
	res, err := s.Store.GetRole(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) CreateRole(ctx context.Context, req *connect.Request[frontendv1.CreateRoleRequest]) (*connect.Response[frontendv1.CreateRoleResponse], error) {
	res, err := s.Store.CreateRole(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateRole(ctx context.Context, req *connect.Request[frontendv1.UpdateRoleRequest]) (*connect.Response[frontendv1.UpdateRoleResponse], error) {
	res, err := s.Store.UpdateRole(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) DeleteRole(ctx context.Context, req *connect.Request[frontendv1.DeleteRoleRequest]) (*connect.Response[frontendv1.DeleteRoleResponse], error) {
	res, err := s.Store.DeleteRole(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
