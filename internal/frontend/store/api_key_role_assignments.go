package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) CreateAPIKeyRoleAssignment(ctx context.Context, req *frontendv1.CreateAPIKeyRoleAssignmentRequest) (*frontendv1.CreateAPIKeyRoleAssignmentResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID := authn.OrganizationID(ctx)

	apiKeyID, err := idformat.APIKey.Parse(req.ApiKeyRoleAssignment.ApiKeyId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	roleID, err := idformat.Role.Parse(req.ApiKeyRoleAssignment.RoleId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	if _, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:             apiKeyID,
		OrganizationID: authn.OrganizationID(ctx),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get api key: %w", err))
		}
		return nil, fmt.Errorf("get api key: %w", err)
	}

	if _, err := q.GetRole(ctx, queries.GetRoleParams{
		ID:             roleID,
		OrganizationID: &orgID,
		ProjectID:      authn.ProjectID(ctx),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role: %w", err))
		}

		return nil, fmt.Errorf("get role: %w", err)
	}

	qAPIKeyRoleAssignment, err := q.CreateAPIKeyRoleAssignment(ctx, queries.CreateAPIKeyRoleAssignmentParams{
		ID:       uuid.New(),
		ApiKeyID: apiKeyID,
		RoleID:   roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("create api key role assignment: %w", err)
	}

	apiKeyRoleAssignment := parseAPIKeyRoleAssignment(qAPIKeyRoleAssignment)

	auditAPIKeyRoleAssignment, err := s.auditlogStore.GetAPIKeyRoleAssignment(ctx, tx, qAPIKeyRoleAssignment.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit log api key role assignment: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_keys.assign_role",
		EventDetails: &auditlogv1.AssignAPIKeyRole{
			ApiKeyRoleAssignment: auditAPIKeyRoleAssignment,
		},
		ResourceType: queries.AuditLogEventResourceTypeApiKey,
		ResourceID:   &qAPIKeyRoleAssignment.ApiKeyID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.CreateAPIKeyRoleAssignmentResponse{
		ApiKeyRoleAssignment: apiKeyRoleAssignment,
	}, nil
}

func (s *Store) DeleteAPIKeyRoleAssignment(ctx context.Context, req *frontendv1.DeleteAPIKeyRoleAssignmentRequest) (*frontendv1.DeleteAPIKeyRoleAssignmentResponse, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyRoleAssignmentID, err := idformat.APIKeyRoleAssignment.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key role assignment id", fmt.Errorf("parse api key role assignment id: %w", err))
	}

	qAPIKeyRoleAssignment, err := q.GetAPIKeyRoleAssignment(ctx, queries.GetAPIKeyRoleAssignmentParams{
		ID:             apiKeyRoleAssignmentID,
		OrganizationID: authn.OrganizationID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key role assignment not found", fmt.Errorf("get api key role assignment: %w", err))
		}
		return nil, fmt.Errorf("get api key role assignment: %w", err)
	}

	auditAPIKeyRoleAssignment, err := s.auditlogStore.GetAPIKeyRoleAssignment(ctx, tx, qAPIKeyRoleAssignment.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit log api key role assignment: %w", err)
	}

	if err := q.DeleteAPIKeyRoleAssignment(ctx, queries.DeleteAPIKeyRoleAssignmentParams{
		ID:             apiKeyRoleAssignmentID,
		OrganizationID: authn.OrganizationID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("delete api key role assignment: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_keys.unassign_role",
		EventDetails: &auditlogv1.UnassignAPIKeyRole{
			ApiKeyRoleAssignment: auditAPIKeyRoleAssignment,
		},
		ResourceType: queries.AuditLogEventResourceTypeApiKey,
		ResourceID:   &qAPIKeyRoleAssignment.ApiKeyID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &frontendv1.DeleteAPIKeyRoleAssignmentResponse{}, nil
}

func (s *Store) ListAPIKeyRoleAssignments(ctx context.Context, req *frontendv1.ListAPIKeyRoleAssignmentsRequest) (*frontendv1.ListAPIKeyRoleAssignmentsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyID, err := idformat.APIKey.Parse(req.ApiKeyId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	var apiKeyRoleAssignments []*frontendv1.APIKeyRoleAssignment
	qAPIKeyRoleAssignments, err := q.ListAPIKeyRoleAssignments(ctx, queries.ListAPIKeyRoleAssignmentsParams{
		ID:             startID,
		ApiKeyID:       apiKeyID,
		OrganizationID: authn.OrganizationID(ctx),
		Limit:          int32(limit + 1),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &frontendv1.ListAPIKeyRoleAssignmentsResponse{
				ApiKeyRoleAssignments: apiKeyRoleAssignments,
			}, nil
		}
		return nil, fmt.Errorf("list api key role assignments: %w", err)
	}

	for _, qAPIKeyRoleAssignment := range qAPIKeyRoleAssignments {
		apiKeyRoleAssignments = append(apiKeyRoleAssignments, parseAPIKeyRoleAssignment(qAPIKeyRoleAssignment))
	}

	var nextPageToken string
	if len(apiKeyRoleAssignments) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(apiKeyRoleAssignments[limit].Id)
		apiKeyRoleAssignments = apiKeyRoleAssignments[:limit]
	}

	return &frontendv1.ListAPIKeyRoleAssignmentsResponse{
		ApiKeyRoleAssignments: apiKeyRoleAssignments,
		NextPageToken:         nextPageToken,
	}, nil
}

func parseAPIKeyRoleAssignment(qAPIKeyRoleAssignment queries.ApiKeyRoleAssignment) *frontendv1.APIKeyRoleAssignment {
	return &frontendv1.APIKeyRoleAssignment{
		Id:       idformat.APIKeyRoleAssignment.Format(qAPIKeyRoleAssignment.ID),
		ApiKeyId: idformat.APIKey.Format(qAPIKeyRoleAssignment.ApiKeyID),
		RoleId:   idformat.Role.Format(qAPIKeyRoleAssignment.RoleID),
	}
}
