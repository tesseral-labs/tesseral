package store

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
)

func (s *Store) GetRBACPolicy(ctx context.Context, req *backendv1.GetRBACPolicyRequest) (*backendv1.GetRBACPolicyResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	return &backendv1.GetRBACPolicyResponse{RbacPolicy: parseRBACPolicy(qActions)}, nil
}

var actionPattern = regexp.MustCompile(`^[a-z0-9_]+\.[a-z0-9_]+\.[a-z0-9_]+`)

func validateActionName(action string) error {
	if !actionPattern.MatchString(action) {
		return apierror.NewInvalidArgumentError("action names must be of the form x.y.z, only containing a-z0-9_", nil)
	}
	if strings.HasPrefix(action, "tesseral") {
		return apierror.NewInvalidArgumentError("action names must not start with 'tesseral'", nil)
	}
	return nil
}

func (s *Store) UpdateRBACPolicy(ctx context.Context, req *backendv1.UpdateRBACPolicyRequest) (*backendv1.UpdateRBACPolicyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	for _, action := range req.RbacPolicy.Actions {
		if err := validateActionName(action.Name); err != nil {
			return nil, fmt.Errorf("validate action name: %w", err)
		}
	}

	var names []string
	for _, action := range req.RbacPolicy.Actions {
		names = append(names, action.Name)
		if err := q.UpsertAction(ctx, queries.UpsertActionParams{
			ID:          uuid.New(),
			ProjectID:   authn.ProjectID(ctx),
			Name:        action.Name,
			Description: action.Description,
		}); err != nil {
			return nil, fmt.Errorf("upsert action: %w", err)
		}
	}

	if _, err := q.DeleteActionsByNameNotInList(ctx, queries.DeleteActionsByNameNotInListParams{
		ProjectID: authn.ProjectID(ctx),
		Names:     names,
	}); err != nil {
		return nil, fmt.Errorf("delete actions by name not in list: %w", err)
	}

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateRBACPolicyResponse{RbacPolicy: parseRBACPolicy(qActions)}, nil
}

func parseRBACPolicy(qActions []queries.Action) *backendv1.RBACPolicy {
	var actions []*backendv1.Action
	for _, qAction := range qActions {
		actions = append(actions, &backendv1.Action{
			Name:        qAction.Name,
			Description: qAction.Description,
		})
	}

	return &backendv1.RBACPolicy{
		Actions: actions,
	}
}
