package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
)

func (s *Service) CreateUserImpersonationToken(ctx context.Context, req *connect.Request[backendv1.CreateUserImpersonationTokenRequest]) (*connect.Response[backendv1.CreateUserImpersonationTokenResponse], error) {
	res, err := s.Store.CreateUserImpersonationToken(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
