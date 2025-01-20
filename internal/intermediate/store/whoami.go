package store

import (
	"context"

	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Store) Whoami(ctx context.Context, req *intermediatev1.WhoamiRequest) (*intermediatev1.WhoamiResponse, error) {
	return &intermediatev1.WhoamiResponse{
		IntermediateSession: authn.IntermediateSession(ctx),
	}, nil
}
