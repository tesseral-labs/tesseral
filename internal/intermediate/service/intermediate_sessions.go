package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Service) SetEmailAsPrimaryLoginFactor(ctx context.Context, req *connect.Request[intermediatev1.SetEmailAsPrimaryLoginFactorRequest]) (*connect.Response[intermediatev1.SetEmailAsPrimaryLoginFactorResponse], error) {
	res, err := s.Store.SetEmailAsPrimaryLoginFactor(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %s", err)
	}

	return connect.NewResponse(res), nil
}

func (s *Service) SetPasswordAsPrimaryLoginFactor(ctx context.Context, req *connect.Request[intermediatev1.SetPasswordAsPrimaryLoginFactorRequest]) (*connect.Response[intermediatev1.SetPasswordAsPrimaryLoginFactorResponse], error) {
	res, err := s.Store.SetPasswordAsPrimaryLoginFactor(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %s", err)
	}

	return connect.NewResponse(res), nil
}
