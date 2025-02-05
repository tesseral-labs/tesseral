package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) SetPrimaryLoginFactor(ctx context.Context, req connect.Request[intermediatev1.SetPrimaryLoginFactorRequest]) (*connect.Response[intermediatev1.SetPrimaryLoginFactorResponse], error) {
	res, err := s.Store.SetPrimaryLoginFactor(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %s", err)
	}

	return connect.NewResponse(res), nil
}
