package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) ListUserRoleAssignments(ctx context.Context, req *connect.Request[frontendv1.ListUserRoleAssignmentsRequest]) (*connect.Response[frontendv1.ListUserRoleAssignmentsResponse], error) {
	res, err := s.Store.ListUserRoleAssignments(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetUserRoleAssignment(ctx context.Context, req *connect.Request[frontendv1.GetUserRoleAssignmentRequest]) (*connect.Response[frontendv1.GetUserRoleAssignmentResponse], error) {
	res, err := s.Store.GetUserRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) CreateUserRoleAssignment(ctx context.Context, req *connect.Request[frontendv1.CreateUserRoleAssignmentRequest]) (*connect.Response[frontendv1.CreateUserRoleAssignmentResponse], error) {
	res, err := s.Store.CreateUserRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) DeleteUserRoleAssignment(ctx context.Context, req *connect.Request[frontendv1.DeleteUserRoleAssignmentRequest]) (*connect.Response[frontendv1.DeleteUserRoleAssignmentResponse], error) {
	res, err := s.Store.DeleteUserRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
