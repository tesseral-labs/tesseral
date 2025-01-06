package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) GetOrganization(ctx context.Context, req *connect.Request[frontendv1.GetOrganizationRequest]) (*connect.Response[frontendv1.GetOrganizationResponse], error) {
	res, err := s.Store.GetOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}

func (s *Service) UpdateOrganization(ctx context.Context, req *connect.Request[frontendv1.UpdateOrganizationRequest]) (*connect.Response[frontendv1.UpdateOrganizationResponse], error) {
	res, err := s.Store.UpdateOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
