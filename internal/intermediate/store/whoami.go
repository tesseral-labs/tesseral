package store

import (
	"context"

	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
)

func (s *Store) Whoami(ctx context.Context, req *intermediatev1.WhoamiRequest) (*intermediatev1.WhoamiResponse, error) {
	return &intermediatev1.WhoamiResponse{
		IntermediateSession: authn.IntermediateSession(ctx),
	}, nil
}
