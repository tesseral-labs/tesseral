package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) GetAPIKeyRoleAssignment(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.APIKeyRoleAssignment, error) {
	qAPIKeyRoleAssignment, err := queries.New(db).GetAPIKeyRoleAssignment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get api key: %w", err)
	}

	return &auditlogv1.APIKeyRoleAssignment{
		Id:       idformat.APIKeyRoleAssignment.Format(qAPIKeyRoleAssignment.ID),
		ApiKeyId: idformat.APIKey.Format(qAPIKeyRoleAssignment.ApiKeyID),
		RoleId:   idformat.Role.Format(qAPIKeyRoleAssignment.RoleID),
	}, nil
}
