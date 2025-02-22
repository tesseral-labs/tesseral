package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) ListSessions(ctx context.Context, req *connect.Request[backendv1.ListSessionsRequest]) (*connect.Response[backendv1.ListSessionsResponse], error) {
	res, err := s.Store.ListSessions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetSession(ctx context.Context, req *connect.Request[backendv1.GetSessionRequest]) (*connect.Response[backendv1.GetSessionResponse], error) {
	res, err := s.Store.GetSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
