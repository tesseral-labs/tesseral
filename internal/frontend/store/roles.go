package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListRoles(ctx context.Context, req *frontendv1.ListRolesRequest) (*frontendv1.ListRolesResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	orgID := authn.OrganizationID(ctx)

	limit := 10
	qRoles, err := q.ListRoles(ctx, queries.ListRolesParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             startID,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	var qRoleIDs []uuid.UUID
	for _, qRole := range qRoles {
		qRoleIDs = append(qRoleIDs, qRole.ID)
	}

	qRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, qRoleIDs)
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role ids: %w", err)
	}

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	var roles []*frontendv1.Role
	for _, qRole := range qRoles {
		roles = append(roles, parseRole(qRole, qRoleActions, qActions))
	}

	var nextPageToken string
	if len(roles) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qRoles[limit].ID)
		roles = roles[:limit]
	}

	return &frontendv1.ListRolesResponse{
		Roles:         roles,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetRole(ctx context.Context, req *frontendv1.GetRoleRequest) (*frontendv1.GetRoleResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	orgID := authn.OrganizationID(ctx)

	qRole, err := q.GetRole(ctx, queries.GetRoleParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             roleID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get role: %w", err)
	}

	qRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{qRole.ID})
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role id: %w", err)
	}

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	return &frontendv1.GetRoleResponse{Role: parseRole(qRole, qRoleActions, qActions)}, nil
}

func (s *Store) CreateRole(ctx context.Context, req *frontendv1.CreateRoleRequest) (*frontendv1.CreateRoleResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qOrg, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	if !qOrg.CustomRolesEnabled {
		return nil, apierror.NewFailedPreconditionError("organization does not have custom roles enabled", nil)
	}

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	var qActionIDs []uuid.UUID
	for _, action := range req.Role.Actions {
		var ok bool
		for _, qAction := range qActions {
			if qAction.Name == action {
				qActionIDs = append(qActionIDs, qAction.ID)
				ok = true
				break
			}
		}
		if !ok {
			return nil, apierror.NewInvalidArgumentError(fmt.Sprintf("invalid action %q", action), fmt.Errorf("action %q not found", action))
		}
	}

	orgID := authn.OrganizationID(ctx)

	qRole, err := q.CreateRole(ctx, queries.CreateRoleParams{
		ID:             uuid.New(),
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		DisplayName:    req.Role.DisplayName,
		Description:    req.Role.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	for _, actionID := range qActionIDs {
		if err := q.UpsertRoleAction(ctx, queries.UpsertRoleActionParams{
			ID:       uuid.New(),
			RoleID:   qRole.ID,
			ActionID: actionID,
		}); err != nil {
			return nil, fmt.Errorf("upsert role action: %w", err)
		}
	}

	qRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{qRole.ID})
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role id: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	pRole := parseRole(qRole, qRoleActions, qActions)
	if _, err := s.CreateTesseralAuditLogEvent(ctx, AuditLogEventData{
		ResourceType: queries.AuditLogEventResourceTypeRole,
		ResourceID:   qRole.ID,
		EventType:    "create",
		Resource:     pRole,
	}); err != nil {
		slog.ErrorContext(ctx, "create_audit_log_event", "error", err)
	}

	return &frontendv1.CreateRoleResponse{Role: pRole}, nil
}

func (s *Store) UpdateRole(ctx context.Context, req *frontendv1.UpdateRoleRequest) (*frontendv1.UpdateRoleResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	qOrg, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	if !qOrg.CustomRolesEnabled {
		return nil, apierror.NewFailedPreconditionError("organization does not have custom roles enabled", nil)
	}

	orgID := authn.OrganizationID(ctx)
	qRole, err := q.GetRoleInOrganization(ctx, queries.GetRoleInOrganizationParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             roleID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role: %w", err))
		}
		return nil, fmt.Errorf("get role: %w", err)
	}

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	qRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{qRole.ID})
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role id: %w", err)
	}

	var updates queries.UpdateRoleParams
	updates.ID = roleID

	updates.DisplayName = qRole.DisplayName
	if req.Role.DisplayName != "" {
		updates.DisplayName = req.Role.DisplayName
	}

	updates.Description = qRole.Description
	if req.Role.Description != "" {
		updates.Description = req.Role.Description
	}

	if req.Role.Actions != nil {
		var qActionIDs []uuid.UUID
		for _, action := range req.Role.Actions {
			var ok bool
			for _, qAction := range qActions {
				if qAction.Name == action {
					qActionIDs = append(qActionIDs, qAction.ID)
					ok = true
					break
				}
			}

			if !ok {
				return nil, apierror.NewInvalidArgumentError(fmt.Sprintf("invalid action %q", action), fmt.Errorf("action %q not found", action))
			}
		}

		for _, qActionID := range qActionIDs {
			if err := q.UpsertRoleAction(ctx, queries.UpsertRoleActionParams{
				ID:       uuid.New(),
				RoleID:   roleID,
				ActionID: qActionID,
			}); err != nil {
				return nil, fmt.Errorf("upsert role action: %w", err)
			}
		}

		if err := q.DeleteRoleActionsByActionIDNotInList(ctx, queries.DeleteRoleActionsByActionIDNotInListParams{
			RoleID:    qRole.ID,
			ActionIds: qActionIDs,
		}); err != nil {
			return nil, fmt.Errorf("delete role actions by action id not in list: %w", err)
		}
	}

	qUpdatedRole, err := q.UpdateRole(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	qUpdatedRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{qUpdatedRole.ID})
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role id: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	pRole := parseRole(qUpdatedRole, qUpdatedRoleActions, qActions)
	pPreviousRole := parseRole(qRole, qRoleActions, qActions)
	if _, err := s.CreateTesseralAuditLogEvent(ctx, AuditLogEventData{
		ResourceType:     queries.AuditLogEventResourceTypeRole,
		ResourceID:       qRole.ID,
		EventType:        "update",
		Resource:         pRole,
		PreviousResource: pPreviousRole,
	}); err != nil {
		slog.ErrorContext(ctx, "create_audit_log_event", "error", err)
	}

	return &frontendv1.UpdateRoleResponse{Role: pRole}, nil
}

func (s *Store) DeleteRole(ctx context.Context, req *frontendv1.DeleteRoleRequest) (*frontendv1.DeleteRoleResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	orgID := authn.OrganizationID(ctx)
	qRole, err := q.GetRoleInOrganization(ctx, queries.GetRoleInOrganizationParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             roleID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role: %w", err))
		}
		return nil, fmt.Errorf("get role: %w", err)
	}

	qRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{qRole.ID})
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role id: %w", err)
	}

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
	}

	err = q.DeleteRole(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("delete role: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	pRole := parseRole(qRole, qRoleActions, qActions)
	if _, err := s.CreateTesseralAuditLogEvent(ctx, AuditLogEventData{
		ResourceType: queries.AuditLogEventResourceTypeRole,
		ResourceID:   qRole.ID,
		EventType:    "delete",
		Resource:     pRole,
	}); err != nil {
		slog.ErrorContext(ctx, "create_audit_log_event", "error", err)
	}

	return &frontendv1.DeleteRoleResponse{}, nil
}

func parseRole(qRole queries.Role, qRoleActions []queries.RoleAction, qActions []queries.Action) *frontendv1.Role {
	var orgID string
	if qRole.OrganizationID != nil {
		orgID = idformat.Organization.Format(*qRole.OrganizationID)
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

	return &frontendv1.Role{
		Id:             idformat.Role.Format(qRole.ID),
		OrganizationId: orgID,
		CreateTime:     timestamppb.New(*qRole.CreateTime),
		UpdateTime:     timestamppb.New(*qRole.UpdateTime),
		DisplayName:    qRole.DisplayName,
		Description:    qRole.Description,
		Actions:        actions,
	}
}
