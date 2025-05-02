package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) GetRBACPolicy(ctx context.Context, req *connect.Request[frontendv1.GetRBACPolicyRequest]) (*connect.Response[frontendv1.GetRBACPolicyResponse], error) {
	res, err := s.Store.GetRBACPolicy(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
