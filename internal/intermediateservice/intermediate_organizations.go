package intermediateservice

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	intermediatev1 "github.com/openauth-dev/openauth/internal/gen/intermediate/v1"
)

func (s *IntermediateService) ListOrganizations(
	ctx context.Context, 
	req *connect.Request[intermediatev1.ListIntermediateOrganizationsRequest],
) (*connect.Response[intermediatev1.ListIntermediateOrganizationsResponse], error) {
	res, err := s.Store.ListIntermediateOrganizations(ctx, req.Msg)
	if err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}

	return connect.NewResponse(res), nil
}