package backendservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
)

func (s *BackendService) CreateOrganization(
	ctx context.Context, 
	req *connect.Request[openauthv1.Organization],
) (*connect.Response[openauthv1.Organization], error) {
		res, err :=  s.Store.CreateOrganization(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}

func (s *BackendService) GetOrganization(
	ctx context.Context, 
	req *connect.Request[openauthv1.ResourceIdRequest],
) (*connect.Response[openauthv1.Organization], error) {
		res, err := s.Store.GetOrganization(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}

func (s *BackendService) ListOrganizations(
	ctx context.Context, 
	req *connect.Request[openauthv1.ListOrganizationsRequest],
) (*connect.Response[openauthv1.ListOrganizationsResponse], error) {
		res, err := s.Store.ListOrganizations(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}

func (s *BackendService) UpdateOrganization(
	ctx context.Context, 
	req *connect.Request[openauthv1.Organization],
) (*connect.Response[openauthv1.Organization], error) {
		res, err := s.Store.UpdateOrganization(ctx, req.Msg)
		if err != nil {
			return nil, fmt.Errorf("store: %w", err)
		}

		return connect.NewResponse(res), nil
}