package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetProjectIDOrganizationBacks(ctx context.Context, organizationID string) (string, error) {
	orgID, err := idformat.Organization.Parse(organizationID)
	if err != nil {
		return "", fmt.Errorf("parse organization id: %w", err)
	}

	orgUUID := uuid.UUID(orgID)
	projectID, err := s.q.GetProjectIDOrganizationBacks(ctx, &orgUUID)
	if err != nil {
		return "", fmt.Errorf("get project id organization backs: %w", err)
	}

	return idformat.Project.Format(projectID), nil
}
