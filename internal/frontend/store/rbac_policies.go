package store

import (
	"context"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
)

func (s *Store) GetRBACPolicy(ctx context.Context, req *frontendv1.GetRBACPolicyRequest) (*frontendv1.GetRBACPolicyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	return &frontendv1.GetRBACPolicyResponse{RbacPolicy: parseRBACPolicy(qActions)}, nil
}

func parseRBACPolicy(qActions []queries.Action) *frontendv1.RBACPolicy {
	var actions []*frontendv1.Action
	for _, qAction := range qActions {
		actions = append(actions, &frontendv1.Action{
			Name:        qAction.Name,
			Description: qAction.Description,
		})
	}

	return &frontendv1.RBACPolicy{
		Actions: actions,
	}
}
