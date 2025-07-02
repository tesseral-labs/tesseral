package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) ExchangeRelayedSessionTokenForSession(ctx context.Context, req *connect.Request[intermediatev1.ExchangeRelayedSessionTokenForSessionRequest]) (*connect.Response[intermediatev1.ExchangeRelayedSessionTokenForSessionResponse], error) {
	res, err := s.Store.ExchangeRelayedSessionTokenForSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	accessToken, err := s.AccessTokenIssuer.NewAccessToken(ctx, authn.ProjectID(ctx), res.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	res.AccessToken = accessToken
	return connect.NewResponse(res), nil
}
