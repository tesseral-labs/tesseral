package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) ListRoles(ctx context.Context, req *backendv1.ListRolesRequest) (*backendv1.ListRolesResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var filterOrgID *uuid.UUID
	if req.OrganizationId != "" {
		orgID, err := idformat.Organization.Parse(req.OrganizationId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
		}

		// not strictly necessary, because orgID is a filtering parameter, not a
		// tenancy boundary, but let's still make sure the organization belongs to
		// this project
		if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ProjectID: authn.ProjectID(ctx),
			ID:        orgID,
		}); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by project id and id: %w", err))
			}
			return nil, fmt.Errorf("get organization by project id and id: %w", err)
		}

		orgUUID := uuid.UUID(orgID)
		filterOrgID = &orgUUID
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qRoles, err := q.ListRoles(ctx, queries.ListRolesParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: filterOrgID,
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

	var roles []*backendv1.Role
	for _, qRole := range qRoles {
		roles = append(roles, parseRole(qRole, qRoleActions, qActions))
	}

	var nextPageToken string
	if len(roles) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qRoles[limit].ID)
		roles = roles[:limit]
	}

	return &backendv1.ListRolesResponse{
		Roles:         roles,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetRole(ctx context.Context, req *backendv1.GetRoleRequest) (*backendv1.GetRoleResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	qRole, err := q.GetRole(ctx, queries.GetRoleParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        roleID,
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

	return &backendv1.GetRoleResponse{Role: parseRole(qRole, qRoleActions, qActions)}, nil
}

func (s *Store) CreateRole(ctx context.Context, req *backendv1.CreateRoleRequest) (*backendv1.CreateRoleResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var roleOrganizationID *uuid.UUID
	if req.Role.OrganizationId != "" {
		orgID, err := idformat.Organization.Parse(req.Role.OrganizationId)
		if err != nil {
			return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
		}

		// unlike for ListRoles, this check is essential; make sure the role's
		// organization belongs to same project as the role itself
		if _, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
			ProjectID: authn.ProjectID(ctx),
			ID:        orgID,
		}); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization: %w", err))
			}
			return nil, fmt.Errorf("get organization: %w", err)
		}

		orgUUID := uuid.UUID(orgID)
		roleOrganizationID = &orgUUID
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

	qRole, err := q.CreateRole(ctx, queries.CreateRoleParams{
		ID:             uuid.New(),
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: roleOrganizationID,
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

	role := parseRole(qRole, qRoleActions, qActions)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.roles.create",
		EventDetails: map[string]any{
			"role": role,
		},
		OrganizationID: roleOrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeRole,
		ResourceID:     &qRole.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.CreateRoleResponse{Role: role}, nil
}

func (s *Store) UpdateRole(ctx context.Context, req *backendv1.UpdateRoleRequest) (*backendv1.UpdateRoleResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	qRole, err := q.GetRole(ctx, queries.GetRoleParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        roleID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role: %w", err))
		}
		return nil, fmt.Errorf("get role: %w", err)
	}

	qPreviousRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{qRole.ID})
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role id: %w", err)
	}

	qActions, err := q.GetActions(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get actions: %w", err)
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

	qRoleActions, err := q.BatchGetRoleActionsByRoleID(ctx, []uuid.UUID{qUpdatedRole.ID})
	if err != nil {
		return nil, fmt.Errorf("batch get role actions by role id: %w", err)
	}

	role := parseRole(qUpdatedRole, qRoleActions, qActions)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.roles.update",
		EventDetails: map[string]any{
			"role":         role,
			"previousRole": parseRole(qRole, qPreviousRoleActions, qActions),
		},
		OrganizationID: qUpdatedRole.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeRole,
		ResourceID:     &qUpdatedRole.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateRoleResponse{Role: role}, nil
}

func (s *Store) DeleteRole(ctx context.Context, req *backendv1.DeleteRoleRequest) (*backendv1.DeleteRoleResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	qRole, err := q.GetRole(ctx, queries.GetRoleParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        roleID,
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

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.roles.delete",
		EventDetails: map[string]any{
			"role": parseRole(qRole, qRoleActions, qActions),
		},
		OrganizationID: qRole.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeRole,
		ResourceID:     &qRole.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteRoleResponse{}, nil
}

func parseRole(qRole queries.Role, qRoleActions []queries.RoleAction, qActions []queries.Action) *backendv1.Role {
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

	return &backendv1.Role{
		Id:             idformat.Role.Format(qRole.ID),
		OrganizationId: orgID,
		CreateTime:     timestamppb.New(*qRole.CreateTime),
		UpdateTime:     timestamppb.New(*qRole.UpdateTime),
		DisplayName:    qRole.DisplayName,
		Description:    qRole.Description,
		Actions:        actions,
	}
}
