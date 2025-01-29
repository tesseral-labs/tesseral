package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) GetOrganizationMicrosoftTenantIDs(ctx context.Context, req *connect.Request[backendv1.GetOrganizationMicrosoftTenantIDsRequest]) (*connect.Response[backendv1.GetOrganizationMicrosoftTenantIDsResponse], error) {
	res, err := s.Store.GetOrganizationMicrosoftTenantIDs(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOrganizationMicrosoftTenantIDs(ctx context.Context, req *connect.Request[backendv1.UpdateOrganizationMicrosoftTenantIDsRequest]) (*connect.Response[backendv1.UpdateOrganizationMicrosoftTenantIDsResponse], error) {
	res, err := s.Store.UpdateOrganizationMicrosoftTenantIDs(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
