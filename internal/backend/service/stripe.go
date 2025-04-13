package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func (s *Service) CreateStripeCheckoutLink(ctx context.Context, req *connect.Request[backendv1.CreateStripeCheckoutLinkRequest]) (*connect.Response[backendv1.CreateStripeCheckoutLinkResponse], error) {
	res, err := s.Store.CreateStripeCheckoutLink(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
