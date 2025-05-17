package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) CreateAPIKeyRoleAssignment(ctx context.Context, req *connect.Request[frontendv1.CreateAPIKeyRoleAssignmentRequest]) (*connect.Response[frontendv1.CreateAPIKeyRoleAssignmentResponse], error) {
	res, err := s.Store.CreateAPIKeyRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) DeleteAPIKeyRoleAssignment(ctx context.Context, req *connect.Request[frontendv1.DeleteAPIKeyRoleAssignmentRequest]) (*connect.Response[frontendv1.DeleteAPIKeyRoleAssignmentResponse], error) {
	res, err := s.Store.DeleteAPIKeyRoleAssignment(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListAPIKeyRoleAssignments(ctx context.Context, req *connect.Request[frontendv1.ListAPIKeyRoleAssignmentsRequest]) (*connect.Response[frontendv1.ListAPIKeyRoleAssignmentsResponse], error) {
	res, err := s.Store.ListAPIKeyRoleAssignments(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
