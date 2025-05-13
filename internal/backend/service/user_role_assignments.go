package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ListUserRoleAssignments(ctx context.Context, req *connect.Request[backendv1.ListUserRoleAssignmentsRequest]) (*connect.Response[backendv1.ListUserRoleAssignmentsResponse], error) {
	res, err := s.Store.ListUserRoleAssignments(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) GetUserRoleAssignment(ctx context.Context, req *connect.Request[backendv1.GetUserRoleAssignmentRequest]) (*connect.Response[backendv1.GetUserRoleAssignmentResponse], error) {
	res, err := s.Store.GetUserRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) CreateUserRoleAssignment(ctx context.Context, req *connect.Request[backendv1.CreateUserRoleAssignmentRequest]) (*connect.Response[backendv1.CreateUserRoleAssignmentResponse], error) {
	res, err := s.Store.CreateUserRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) DeleteUserRoleAssignment(ctx context.Context, req *connect.Request[backendv1.DeleteUserRoleAssignmentRequest]) (*connect.Response[backendv1.DeleteUserRoleAssignmentResponse], error) {
	res, err := s.Store.DeleteUserRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
