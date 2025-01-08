package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/projectid"
)

func (s *Service) SignInWithEmail(ctx context.Context, req *connect.Request[intermediatev1.SignInWithEmailRequest]) (*connect.Response[intermediatev1.SignInWithEmailResponse], error) {
	res, err := s.Store.SignInWithEmail(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	projectID := projectid.ProjectID(ctx)

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.BuildCookie(projectID, "intermediateAccessToken", res.IntermediateSessionToken, req.Spec().Schema == "https"))

	return connectResponse, nil
}
