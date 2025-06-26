package store

import (
	"context"
	"fmt"

	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
)

func (s *Store) UpdateMe(ctx context.Context, req *frontendv1.UpdateMeRequest) (*frontendv1.UpdateMeResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID := authn.UserID(ctx)

	auditPreviousUser, err := s.auditlogStore.GetUser(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get audit previous user: %w", err)
	}

	updates := queries.UpdateMeParams{
		ID: userID,
	}

	if req.User.DisplayName != nil && derefOrEmpty(req.User.DisplayName) != "" {
		updates.DisplayName = req.User.DisplayName
	}

	qUpdatedUser, err := q.UpdateMe(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	auditUser, err := s.auditlogStore.GetUser(ctx, tx, qUpdatedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit user: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.update",
		EventDetails: &auditlogv1.UpdateUser{
			User:         auditUser,
			PreviousUser: auditPreviousUser,
		},
		ResourceType: queries.AuditLogEventResourceTypeUser,
		ResourceID:   &userID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	// send sync user event
	if err := s.sendSyncUserEvent(ctx, qUpdatedUser); err != nil {
		return nil, fmt.Errorf("send sync user event: %w", err)
	}

	return &frontendv1.UpdateMeResponse{
		User: parseUser(qUpdatedUser),
	}, nil
}
