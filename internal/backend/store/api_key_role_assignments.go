package store

import (
	"context"

	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Store) CreateAPIKeyRoleAssignment(ctx context.Context, req *backendv1.CreateAPIKeyRoleAssignmentRequest) (*backendv1.CreateAPIKeyRoleAssignmentResponse, error) {
	return &backendv1.CreateAPIKeyRoleAssignmentResponse{}, nil
}

func (s *Store) DeleteAPIKeyRoleAssignment(ctx context.Context, req *backendv1.DeleteAPIKeyRoleAssignmentRequest) (*backendv1.DeleteAPIKeyRoleAssignmentResponse, error) {
	return &backendv1.DeleteAPIKeyRoleAssignmentResponse{}, nil
}

func (s *Store) ListAPIKeyRoleAssignments(ctx context.Context, req *backendv1.ListAPIKeyRoleAssignmentsRequest) (*backendv1.ListAPIKeyRoleAssignmentsResponse, error) {
	return &backendv1.ListAPIKeyRoleAssignmentsResponse{}, nil
}
