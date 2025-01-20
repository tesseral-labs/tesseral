package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) CreateIntermediateSession(ctx context.Context, req *connect.Request[intermediatev1.CreateIntermediateSessionRequest]) (*connect.Response[intermediatev1.CreateIntermediateSessionResponse], error) {
	res, err := s.Store.CreateIntermediateSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.BuildCookie(ctx, req, "intermediateAccessToken", res.IntermediateSessionSecretToken, authn.ProjectID(ctx)))

	return connectResponse, nil
}
