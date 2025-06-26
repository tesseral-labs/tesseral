package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) GetUserRoleAssignment(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.UserRoleAssignment, error) {
	qUserRoleAssignment, err := queries.New(db).GetUserRoleAssignment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user role assignment: %w", err)
	}

	return &auditlogv1.UserRoleAssignment{
		Id:     idformat.UserRoleAssignment.Format(qUserRoleAssignment.ID),
		UserId: idformat.User.Format(qUserRoleAssignment.UserID),
		RoleId: idformat.Role.Format(qUserRoleAssignment.RoleID),
	}, nil
}
