package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetRole(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.Role, error) {
	qRole, err := queries.New(db).GetRole(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}

	qRoleActions, err := queries.New(db).BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{id})
	if err != nil {
		return nil, fmt.Errorf("get role actions: %w", err)
	}

	qActions, err := queries.New(db).GetActions(ctx, qRole.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	var actions []string
	for _, qRoleAction := range qRoleActions {
		if qRoleAction.RoleID != qRole.ID {
			continue
		}

		for _, qAction := range qActions {
			if qAction.ID == qRoleAction.ActionID {
				actions = append(actions, qAction.Name)
				break
			}
		}
	}

	return &auditlogv1.Role{
		Id:          idformat.Role.Format(qRole.ID),
		CreateTime:  timestamppb.New(*qRole.CreateTime),
		UpdateTime:  timestamppb.New(*qRole.UpdateTime),
		DisplayName: qRole.DisplayName,
		Description: qRole.Description,
		Actions:     actions,
	}, nil
}
