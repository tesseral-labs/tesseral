package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) GetOrganizationDomains(ctx context.Context, req *connect.Request[backendv1.GetOrganizationDomainsRequest]) (*connect.Response[backendv1.GetOrganizationDomainsResponse], error) {
	res, err := s.Store.GetOrganizationDomains(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOrganizationDomains(ctx context.Context, req *connect.Request[backendv1.UpdateOrganizationDomainsRequest]) (*connect.Response[backendv1.UpdateOrganizationDomainsResponse], error) {
	res, err := s.Store.UpdateOrganizationDomains(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
