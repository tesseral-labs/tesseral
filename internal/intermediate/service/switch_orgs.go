package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) ExchangeSessionForIntermediateSession(ctx context.Context, req *connect.Request[intermediatev1.ExchangeSessionForIntermediateSessionRequest]) (*connect.Response[intermediatev1.ExchangeSessionForIntermediateSessionResponse], error) {
	refreshToken, _ := cookies.GetRefreshToken(authn.ProjectID(ctx), req)
	if refreshToken != "" {
		req.Msg.RefreshToken = refreshToken
	}

	if refreshToken == "" {
		return nil, apierror.NewUnauthenticatedError("no refresh token provided", nil)
	}

	res, err := s.Store.ExchangeSessionForIntermediateSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.NewIntermediateAccessToken(authn.ProjectID(ctx), res.IntermediateSessionSecretToken))

	return connectResponse, nil
}
