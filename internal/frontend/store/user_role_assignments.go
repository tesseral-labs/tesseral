package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) ListUserRoleAssignments(ctx context.Context, req *frontendv1.ListUserRoleAssignmentsRequest) (*frontendv1.ListUserRoleAssignmentsResponse, error) {
	if req.RoleId != "" {
		return s.listUserRoleAssignmentsByRoleID(ctx, req)
	} else if req.UserId != "" {
		return s.listUserRoleAssignmentsByUserID(ctx, req)
	} else {
		return nil, apierror.NewInvalidArgumentError("one of role_id or user_id must be provided", nil)
	}
}

func (s *Store) listUserRoleAssignmentsByRoleID(ctx context.Context, req *frontendv1.ListUserRoleAssignmentsRequest) (*frontendv1.ListUserRoleAssignmentsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.RoleId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	// authz
	orgID := authn.OrganizationID(ctx)
	if _, err := q.GetRoleInOrganization(ctx, queries.GetRoleInOrganizationParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             roleID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get role: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qUserRoleAssignments, err := q.ListUserRoleAssignmentsByRole(ctx, queries.ListUserRoleAssignmentsByRoleParams{
		RoleID: roleID,
		ID:     startID,
		Limit:  int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list user role assignments: %w", err)
	}

	var userRoleAssignments []*frontendv1.UserRoleAssignment
	for _, qUserRoleAssignment := range qUserRoleAssignments {
		userRoleAssignments = append(userRoleAssignments, parseUserRoleAssignment(qUserRoleAssignment))
	}

	var nextPageToken string
	if len(userRoleAssignments) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qUserRoleAssignments[limit].ID)
		userRoleAssignments = userRoleAssignments[:limit]
	}

	return &frontendv1.ListUserRoleAssignmentsResponse{
		UserRoleAssignments: userRoleAssignments,
		NextPageToken:       nextPageToken,
	}, nil
}

func (s *Store) listUserRoleAssignmentsByUserID(ctx context.Context, req *frontendv1.ListUserRoleAssignmentsRequest) (*frontendv1.ListUserRoleAssignmentsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.UserId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	// authz
	if _, err := q.GetUser(ctx, queries.GetUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user by project id and id: %w", err))
		}

		return nil, fmt.Errorf("get user: %w", err)
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, fmt.Errorf("unmarshal page token: %w", err)
	}

	limit := 10
	qUserRoleAssignments, err := q.ListUserRoleAssignmentsByUser(ctx, queries.ListUserRoleAssignmentsByUserParams{
		UserID: userID,
		ID:     startID,
		Limit:  int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list user role assignments: %w", err)
	}

	var userRoleAssignments []*frontendv1.UserRoleAssignment
	for _, qUserRoleAssignment := range qUserRoleAssignments {
		userRoleAssignments = append(userRoleAssignments, parseUserRoleAssignment(qUserRoleAssignment))
	}

	var nextPageToken string
	if len(userRoleAssignments) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qUserRoleAssignments[limit].ID)
		userRoleAssignments = userRoleAssignments[:limit]
	}

	return &frontendv1.ListUserRoleAssignmentsResponse{
		UserRoleAssignments: userRoleAssignments,
		NextPageToken:       nextPageToken,
	}, nil
}

func (s *Store) GetUserRoleAssignment(ctx context.Context, req *frontendv1.GetUserRoleAssignmentRequest) (*frontendv1.GetUserRoleAssignmentResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	id, err := idformat.UserRoleAssignment.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user role assignment id", fmt.Errorf("parse user role assignment id: %w", err))
	}

	orgID := authn.OrganizationID(ctx)
	qUserRoleAssignment, err := q.GetUserRoleAssignment(ctx, queries.GetUserRoleAssignmentParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user role assignment not found", fmt.Errorf("get user role assignment: %w", err))
		}
		return nil, fmt.Errorf("get user role assignment: %w", err)
	}

	return &frontendv1.GetUserRoleAssignmentResponse{UserRoleAssignment: parseUserRoleAssignment(qUserRoleAssignment)}, nil
}

func (s *Store) CreateUserRoleAssignment(ctx context.Context, req *frontendv1.CreateUserRoleAssignmentRequest) (*frontendv1.CreateUserRoleAssignmentResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	roleID, err := idformat.Role.Parse(req.UserRoleAssignment.RoleId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	userID, err := idformat.User.Parse(req.UserRoleAssignment.UserId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user id", fmt.Errorf("parse user id: %w", err))
	}

	// ensure both role and user belong to project/organization
	orgID := authn.OrganizationID(ctx)
	if _, err := q.GetRole(ctx, queries.GetRoleParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             roleID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role: %w", err))
		}
		return nil, fmt.Errorf("get role: %w", err)
	}

	if _, err := q.GetUser(ctx, queries.GetUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user not found", fmt.Errorf("get user: %w", err))
		}
		return nil, err
	}

	if err := q.UpsertUserRoleAssignment(ctx, queries.UpsertUserRoleAssignmentParams{
		ID:     uuid.New(),
		RoleID: roleID,
		UserID: userID,
	}); err != nil {
		return nil, fmt.Errorf("upsert user role assignment: %w", err)
	}

	qUserRoleAssignment, err := q.GetUserRoleAssignmentByUserAndRole(ctx, queries.GetUserRoleAssignmentByUserAndRoleParams{
		UserID: userID,
		RoleID: roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("get user role assignment by user and role: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.CreateUserRoleAssignmentResponse{UserRoleAssignment: parseUserRoleAssignment(qUserRoleAssignment)}, nil
}

func (s *Store) DeleteUserRoleAssignment(ctx context.Context, req *frontendv1.DeleteUserRoleAssignmentRequest) (*frontendv1.DeleteUserRoleAssignmentResponse, error) {
	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	id, err := idformat.UserRoleAssignment.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid user role assignment id", fmt.Errorf("parse user role assignment id: %w", err))
	}

	orgID := authn.OrganizationID(ctx)
	if _, err := q.GetUserRoleAssignment(ctx, queries.GetUserRoleAssignmentParams{
		ProjectID:      authn.ProjectID(ctx),
		OrganizationID: &orgID,
		ID:             id,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("user role assignment not found", fmt.Errorf("get user role assignment: %w", err))
		}
		return nil, fmt.Errorf("get user role assignment: %w", err)
	}

	if err := q.DeleteUserRoleAssignment(ctx, id); err != nil {
		return nil, fmt.Errorf("delete user role assignment: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &frontendv1.DeleteUserRoleAssignmentResponse{}, nil
}

func parseUserRoleAssignment(qUserRoleAssignment queries.UserRoleAssignment) *frontendv1.UserRoleAssignment {
	return &frontendv1.UserRoleAssignment{
		Id:     idformat.UserRoleAssignment.Format(qUserRoleAssignment.ID),
		RoleId: idformat.Role.Format(qUserRoleAssignment.RoleID),
		UserId: idformat.User.Format(qUserRoleAssignment.UserID),
	}
}
