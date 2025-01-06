package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) GetProject(ctx context.Context, req *connect.Request[frontendv1.GetProjectRequest]) (*connect.Response[frontendv1.GetProjectResponse], error) {
	res, err := s.Store.GetProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
