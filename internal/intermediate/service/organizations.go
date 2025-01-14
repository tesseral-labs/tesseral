package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) ListOrganizations(ctx context.Context, req *connect.Request[intermediatev1.ListOrganizationsRequest]) (*connect.Response[intermediatev1.ListOrganizationsResponse], error) {
	res, err := s.Store.ListOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) ListSAMLOrganizations(ctx context.Context, req *connect.Request[intermediatev1.ListSAMLOrganizationsRequest]) (*connect.Response[intermediatev1.ListSAMLOrganizationsResponse], error) {
	res, err := s.Store.ListSAMLOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
