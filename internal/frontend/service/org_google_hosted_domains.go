package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) GetOrganizationGoogleHostedDomains(ctx context.Context, req *connect.Request[frontendv1.GetOrganizationGoogleHostedDomainsRequest]) (*connect.Response[frontendv1.GetOrganizationGoogleHostedDomainsResponse], error) {
	res, err := s.Store.GetOrganizationGoogleHostedDomains(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOrganizationGoogleHostedDomains(ctx context.Context, req *connect.Request[frontendv1.UpdateOrganizationGoogleHostedDomainsRequest]) (*connect.Response[frontendv1.UpdateOrganizationGoogleHostedDomainsResponse], error) {
	res, err := s.Store.UpdateOrganizationGoogleHostedDomains(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
