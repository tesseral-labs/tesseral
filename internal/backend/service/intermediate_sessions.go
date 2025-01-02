package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) ListIntermediateSessions(ctx context.Context, req *connect.Request[backendv1.ListIntermediateSessionsRequest]) (*connect.Response[backendv1.ListIntermediateSessionsResponse], error) {
	res, err := s.Store.ListIntermediateSessions(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) GetIntermediateSession(ctx context.Context, req *connect.Request[backendv1.GetIntermediateSessionRequest]) (*connect.Response[backendv1.GetIntermediateSessionResponse], error) {
	res, err := s.Store.GetIntermediateSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
