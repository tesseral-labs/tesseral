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
	"github.com/tesseral-labs/tesseral/internal/muststructpb"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) CreateAPIKeyRoleAssignment(ctx context.Context, req *backendv1.CreateAPIKeyRoleAssignmentRequest) (*backendv1.CreateAPIKeyRoleAssignmentResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyID, err := idformat.APIKey.Parse(req.ApiKeyRoleAssignment.ApiKeyId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key id", fmt.Errorf("parse api key id: %w", err))
	}

	roleID, err := idformat.Role.Parse(req.ApiKeyRoleAssignment.RoleId)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid role id", fmt.Errorf("parse role id: %w", err))
	}

	qAPIKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:        apiKeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get api key: %w", err))
		}
		return nil, fmt.Errorf("get api key: %w", err)
	}

	if _, err := q.GetRole(ctx, queries.GetRoleParams{
		ID:        roleID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("role not found", fmt.Errorf("get role: %w", err))
		}

		return nil, fmt.Errorf("get role: %w", err)
	}

	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ID:        qAPIKey.OrganizationID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get organization: %w", err))
		}
		return nil, fmt.Errorf("get organization: %w", err)
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
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_key_role_assignments.create",
		EventDetails: muststructpb.MustNewValue(map[string]any{
			"apiKeyRoleAssignment": apiKeyRoleAssignment,
		}),
		OrganizationID: &qOrg.ID,
		ResourceType:   queries.AuditLogEventResourceTypeApiKeyRoleAssignment,
		ResourceID:     &qAPIKeyRoleAssignment.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.CreateAPIKeyRoleAssignmentResponse{
		ApiKeyRoleAssignment: apiKeyRoleAssignment,
	}, nil
}

func (s *Store) DeleteAPIKeyRoleAssignment(ctx context.Context, req *backendv1.DeleteAPIKeyRoleAssignmentRequest) (*backendv1.DeleteAPIKeyRoleAssignmentResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	apiKeyRoleAssignmentID, err := idformat.APIKeyRoleAssignment.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid api key role assignment id", fmt.Errorf("parse api key role assignment id: %w", err))
	}

	qAPIKeyRoleAssignment, err := q.GetAPIKeyRoleAssignment(ctx, queries.GetAPIKeyRoleAssignmentParams{
		ID:        apiKeyRoleAssignmentID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apierror.NewNotFoundError("api key role assignment not found", fmt.Errorf("get api key role assignment: %w", err))
		}
		return nil, fmt.Errorf("get api key role assignment: %w", err)
	}

	qAPIKey, err := q.GetAPIKeyByID(ctx, queries.GetAPIKeyByIDParams{
		ID:        qAPIKeyRoleAssignment.ApiKeyID,
		ProjectID: authn.ProjectID(ctx),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("api key not found", fmt.Errorf("get api key: %w", err))
		}
		return nil, fmt.Errorf("get api key: %w", err)
	}

	if err := q.DeleteAPIKeyRoleAssignment(ctx, queries.DeleteAPIKeyRoleAssignmentParams{
		ID:        apiKeyRoleAssignmentID,
		ProjectID: authn.ProjectID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("delete api key role assignment: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.api_key_role_assignments.delete",
		EventDetails: muststructpb.MustNewValue(map[string]any{
			"apiKeyRoleAssignment": parseAPIKeyRoleAssignment(qAPIKeyRoleAssignment),
		}),
		OrganizationID: &qAPIKey.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeApiKeyRoleAssignment,
		ResourceID:     &qAPIKeyRoleAssignment.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &backendv1.DeleteAPIKeyRoleAssignmentResponse{}, nil
}

func (s *Store) ListAPIKeyRoleAssignments(ctx context.Context, req *backendv1.ListAPIKeyRoleAssignmentsRequest) (*backendv1.ListAPIKeyRoleAssignmentsResponse, error) {
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
	qAPIKeyRoleAssignments, err := q.ListAPIKeyRoleAssignments(ctx, queries.ListAPIKeyRoleAssignmentsParams{
		ID:        startID,
		ApiKeyID:  apiKeyID,
		ProjectID: authn.ProjectID(ctx),
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list api key role assignments: %w", err)
	}

	var apiKeyRoleAssignments []*backendv1.APIKeyRoleAssignment
	for _, qAPIKeyRoleAssignment := range qAPIKeyRoleAssignments {
		apiKeyRoleAssignments = append(apiKeyRoleAssignments, parseAPIKeyRoleAssignment(qAPIKeyRoleAssignment))
	}

	var nextPageToken string
	if len(apiKeyRoleAssignments) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(apiKeyRoleAssignments[limit].Id)
		apiKeyRoleAssignments = apiKeyRoleAssignments[:limit]
	}

	return &backendv1.ListAPIKeyRoleAssignmentsResponse{
		ApiKeyRoleAssignments: apiKeyRoleAssignments,
		NextPageToken:         nextPageToken,
	}, nil
}

func parseAPIKeyRoleAssignment(qAPIKeyRoleAssignment queries.ApiKeyRoleAssignment) *backendv1.APIKeyRoleAssignment {
	return &backendv1.APIKeyRoleAssignment{
		Id:       idformat.APIKeyRoleAssignment.Format(qAPIKeyRoleAssignment.ID),
		ApiKeyId: idformat.APIKey.Format(qAPIKeyRoleAssignment.ApiKeyID),
		RoleId:   idformat.Role.Format(qAPIKeyRoleAssignment.RoleID),
	}
}
