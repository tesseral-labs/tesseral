package frontendservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	frontendv1 "github.com/openauth/openauth/internal/gen/frontend/v1"
)

func (s *FrontendService) GetAccessToken(ctx context.Context, req *connect.Request[frontendv1.GetAccessTokenRequest]) (*connect.Response[frontendv1.GetAccessTokenResponse], error) {
	res, err := s.Store.GetAccessToken(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return connect.NewResponse(res), nil
}
