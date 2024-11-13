package store

import (
	"context"

	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
)

func (s *Store) CreateOrganization(ctx context.Context, req *openauthv1.Organization) (*openauthv1.Organization, error) {
	return nil, nil
}

func (s *Store) GetOrganization(ctx context.Context, req *openauthv1.ResourceIdRequest) (*openauthv1.Organization, error) {
	return nil, nil
}

// TODO: Ensure that this function can only be called via a backend service reuqest
func (s *Store) ListOrganizations(ctx context.Context, req *openauthv1.ProjectIdRequest) (*backendv1.ListOrganizationsResponse, error) {
	return nil, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *openauthv1.Organization) (*openauthv1.Organization, error) {
	return nil, nil
}