package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
)

func (s *Service) SetUserPassword(ctx context.Context, req *connect.Request[frontendv1.SetUserPasswordRequest]) (*connect.Response[frontendv1.SetUserPasswordResponse], error) {
	res, err := s.Store.SetUserPassword(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}
