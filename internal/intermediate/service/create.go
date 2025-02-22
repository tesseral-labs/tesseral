package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) CreateIntermediateSession(ctx context.Context, req *connect.Request[intermediatev1.CreateIntermediateSessionRequest]) (*connect.Response[intermediatev1.CreateIntermediateSessionResponse], error) {
	res, err := s.Store.CreateIntermediateSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.NewIntermediateAccessToken(authn.ProjectID(ctx), res.IntermediateSessionSecretToken))

	return connectResponse, nil
}
