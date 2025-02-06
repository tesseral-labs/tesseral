package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Service) SetEmailAsPrimaryLoginFactor(ctx context.Context, req *connect.Request[intermediatev1.SetEmailAsPrimaryLoginFactorRequest]) (*connect.Response[intermediatev1.SetEmailAsPrimaryLoginFactorResponse], error) {
	res, err := s.Store.SetEmailAsPrimaryLoginFactor(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %s", err)
	}

	return connect.NewResponse(res), nil
}
