package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) CreateAPIKeyRoleAssignment(ctx context.Context, req *connect.Request[backendv1.CreateAPIKeyRoleAssignmentRequest]) (*connect.Response[backendv1.CreateAPIKeyRoleAssignmentResponse], error) {
	res, err := s.Store.CreateAPIKeyRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteAPIKeyRoleAssignment(ctx context.Context, req *connect.Request[backendv1.DeleteAPIKeyRoleAssignmentRequest]) (*connect.Response[backendv1.DeleteAPIKeyRoleAssignmentResponse], error) {
	res, err := s.Store.DeleteAPIKeyRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListAPIKeyRoleAssignments(ctx context.Context, req *connect.Request[backendv1.ListAPIKeyRoleAssignmentsRequest]) (*connect.Response[backendv1.ListAPIKeyRoleAssignmentsResponse], error) {
	res, err := s.Store.ListAPIKeyRoleAssignments(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
