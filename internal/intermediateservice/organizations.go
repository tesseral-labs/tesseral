package intermediateservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth-dev/openauth/internal/gen/intermediate/v1"
)

func (s *IntermediateService) CreateOrganization(
	ctx context.Context,
	req *connect.Request[intermediatev1.CreateOrganizationRequest],
) (*connect.Response[intermediatev1.CreateOrganizationResponse], error) {
	res, err := s.Store.CreateIntermediateOrganization(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(&intermediatev1.CreateOrganizationResponse{
		Organization: res,
	}), nil
}

func (s *IntermediateService) ListOrganizations(
	ctx context.Context, 
	req *connect.Request[intermediatev1.ListOrganizationsRequest],
) (*connect.Response[intermediatev1.ListOrganizationsResponse], error) {
	res, err := s.Store.ListIntermediateOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}