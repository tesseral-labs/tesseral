package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ConsoleSearch(ctx context.Context, req *connect.Request[backendv1.ConsoleSearchRequest]) (*connect.Response[backendv1.ConsoleSearchResponse], error) {
	resp, err := s.Store.ConsoleSearch(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(resp), nil
}
