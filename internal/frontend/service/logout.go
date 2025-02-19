package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
)

func (s *Service) Logout(ctx context.Context, req *connect.Request[frontendv1.LogoutRequest]) (*connect.Response[frontendv1.LogoutResponse], error) {
	res, err := s.Store.Logout(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectRes := connect.NewResponse(res)

	connectRes.Header().Add("Set-Cookie", cookies.ExpiredAccessToken(authn.ProjectID(ctx)))
	connectRes.Header().Add("Set-Cookie", cookies.ExpiredIntermediateAccessToken(authn.ProjectID(ctx)))
	connectRes.Header().Add("Set-Cookie", cookies.ExpiredRefreshToken(authn.ProjectID(ctx)))

	return connectRes, nil
}
