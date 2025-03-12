package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) ListSwitchableOrganizations(ctx context.Context, req *connect.Request[frontendv1.ListSwitchableOrganizationsRequest]) (*connect.Response[frontendv1.ListSwitchableOrganizationsResponse], error) {
	res, err := s.Store.ListSwitchableOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
