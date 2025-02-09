package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) CreateProject(ctx context.Context, req *connect.Request[intermediatev1.CreateProjectRequest]) (*connect.Response[intermediatev1.CreateProjectResponse], error) {
	res, err := s.Store.CreateProject(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
