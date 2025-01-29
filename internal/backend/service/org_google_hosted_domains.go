package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) GetOrganizationGoogleHostedDomains(ctx context.Context, req *connect.Request[backendv1.GetOrganizationGoogleHostedDomainsRequest]) (*connect.Response[backendv1.GetOrganizationGoogleHostedDomainsResponse], error) {
	res, err := s.Store.GetOrganizationGoogleHostedDomains(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOrganizationGoogleHostedDomains(ctx context.Context, req *connect.Request[backendv1.UpdateOrganizationGoogleHostedDomainsRequest]) (*connect.Response[backendv1.UpdateOrganizationGoogleHostedDomainsResponse], error) {
	res, err := s.Store.UpdateOrganizationGoogleHostedDomains(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
