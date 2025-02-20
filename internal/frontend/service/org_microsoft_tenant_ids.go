package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) GetOrganizationMicrosoftTenantIDs(ctx context.Context, req *connect.Request[frontendv1.GetOrganizationMicrosoftTenantIDsRequest]) (*connect.Response[frontendv1.GetOrganizationMicrosoftTenantIDsResponse], error) {
	res, err := s.Store.GetOrganizationMicrosoftTenantIDs(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOrganizationMicrosoftTenantIDs(ctx context.Context, req *connect.Request[frontendv1.UpdateOrganizationMicrosoftTenantIDsRequest]) (*connect.Response[frontendv1.UpdateOrganizationMicrosoftTenantIDsResponse], error) {
	res, err := s.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
