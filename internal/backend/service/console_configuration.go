package service

import (
	"context"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ConsoleGetConfiguration(ctx context.Context, req *connect.Request[backendv1.ConsoleGetConfigurationRequest]) (*connect.Response[backendv1.ConsoleGetConfigurationResponse], error) {
	resp, err := s.Store.ConsoleGetConfiguration(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}
