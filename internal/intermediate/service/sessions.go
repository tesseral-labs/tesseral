package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/cookies"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/projectid"
)

func (s *Service) ExchangeIntermediateSessionForNewOrganizationSession(ctx context.Context, req *connect.Request[intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionRequest]) (*connect.Response[intermediatev1.ExchangeIntermediateSessionForNewOrganizationSessionResponse], error) {
	res, err := s.Store.ExchangeIntermediateSessionForNewOrganizationSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	projectID := projectid.ProjectID(ctx)

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.BuildCookie(projectID, "accessToken", res.AccessToken, req.Spec().Schema == "https"))

	return connectResponse, nil
}

func (s *Service) ExchangeIntermediateSessionForSession(ctx context.Context, req *connect.Request[intermediatev1.ExchangeIntermediateSessionForSessionRequest]) (*connect.Response[intermediatev1.ExchangeIntermediateSessionForSessionResponse], error) {
	res, err := s.Store.ExchangeIntermediateSessionForSession(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	projectID := projectid.ProjectID(ctx)

	connectResponse := connect.NewResponse(res)
	connectResponse.Header().Add("Set-Cookie", cookies.BuildCookie(projectID, "accessToken", res.AccessToken, req.Spec().Schema == "https"))

	return connectResponse, nil
}
