package store

import (
	"context"

	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
)

func (s * Store) CreateUser(ctx context.Context, req *openauthv1.User) (*openauthv1.User, error) {
	return nil, nil
}

func (s * Store) GetUser(ctx context.Context, req *openauthv1.ResourceIdRequest) (*openauthv1.User, error) {
	return nil, nil
}

// TODO: Ensure that this function can only be called via a backend service reuqest
func (s * Store) ListUsers(ctx context.Context, req *openauthv1.ProjectIdRequest) (*backendv1.ListUsersResponse, error) {
	return nil, nil
}

func (s * Store) UpdateUser(ctx context.Context, req *openauthv1.User) (*openauthv1.User, error) {
	return nil, nil
}