package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateOrganization(ctx context.Context, req *backendv1.CreateOrganizationRequest) (*backendv1.CreateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// Creating organizations directly on the Dogfood Project
	// without a project to back with the organization will
	// break the intermediate service, so we restrict the
	// ability to create an organization in this case.
	if authn.ProjectID(ctx) == *s.dogfoodProjectID {
		return nil, apierror.NewPermissionDeniedError("dogfood project cannot create organizations", fmt.Errorf("dogfood project cannot create organizations"))
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	if derefOrEmpty(req.Organization.LogInWithGoogle) && !qProject.LogInWithGoogle {
		return nil, apierror.NewPermissionDeniedError("log in with google is not enabled for this project", fmt.Errorf("log in with google is not enabled for this project"))
	}

	if derefOrEmpty(req.Organization.LogInWithMicrosoft) && !qProject.LogInWithMicrosoft {
		return nil, apierror.NewPermissionDeniedError("log in with microsoft is not enabled for this project", fmt.Errorf("log in with microsoft is not enabled for this project"))
	}

	if derefOrEmpty(req.Organization.LogInWithEmail) && !qProject.LogInWithEmail {
		return nil, apierror.NewPermissionDeniedError("log in with email is not enabled for this project", fmt.Errorf("log in with email is not enabled for this project"))
	}

	if derefOrEmpty(req.Organization.LogInWithPassword) && !qProject.LogInWithPassword {
		return nil, apierror.NewPermissionDeniedError("log in with password is not enabled for this project", fmt.Errorf("log in with password is not enabled for this project"))
	}

	if derefOrEmpty(req.Organization.LogInWithSaml) && !qProject.LogInWithSaml {
		return nil, apierror.NewPermissionDeniedError("log in with saml is not enabled for this project", fmt.Errorf("log in with saml is not enabled for this project"))
	}

	if derefOrEmpty(req.Organization.LogInWithAuthenticatorApp) && !qProject.LogInWithAuthenticatorApp {
		return nil, apierror.NewPermissionDeniedError("log in with authenticator app is not enabled for this project", fmt.Errorf("log in with authenticator app is not enabled for this project"))
	}

	if derefOrEmpty(req.Organization.LogInWithPasskey) && !qProject.LogInWithPasskey {
		return nil, apierror.NewPermissionDeniedError("log in with passkey is not enabled for this project", fmt.Errorf("log in with passkey is not enabled for this project"))
	}

	var scimEnabled bool
	if req.Organization.ScimEnabled != nil {
		scimEnabled = *req.Organization.ScimEnabled
	}

	qOrg, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                        uuid.New(),
		ProjectID:                 authn.ProjectID(ctx),
		DisplayName:               req.Organization.DisplayName,
		LogInWithGoogle:           derefOrEmpty(req.Organization.LogInWithGoogle),
		LogInWithMicrosoft:        derefOrEmpty(req.Organization.LogInWithMicrosoft),
		LogInWithEmail:            derefOrEmpty(req.Organization.LogInWithEmail),
		LogInWithPassword:         derefOrEmpty(req.Organization.LogInWithPassword),
		LogInWithSaml:             derefOrEmpty(req.Organization.LogInWithSaml),
		LogInWithAuthenticatorApp: derefOrEmpty(req.Organization.LogInWithAuthenticatorApp),
		LogInWithPasskey:          derefOrEmpty(req.Organization.LogInWithPasskey),
		ScimEnabled:               scimEnabled,
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	// Send webhook event
	if err := s.sendSyncOrganizationEvent(ctx, qOrg); err != nil {
		return nil, fmt.Errorf("send sync organization event: %w", err)
	}

	return &backendv1.CreateOrganizationResponse{Organization: parseOrganization(qProject, qOrg)}, nil
}

func (s *Store) ListOrganizations(ctx context.Context, req *backendv1.ListOrganizationsRequest) (*backendv1.ListOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	limit := 10
	qOrgs, err := q.ListOrganizationsByProjectId(ctx, queries.ListOrganizationsByProjectIdParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        startID,
		Limit:     int32(limit + 1),
	})
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}

	var organizations []*backendv1.Organization
	for _, qOrg := range qOrgs {
		organizations = append(organizations, parseOrganization(qProject, qOrg))
	}

	var nextPageToken string
	if len(organizations) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(qOrgs[limit].ID)
		organizations = organizations[:limit]
	}

	return &backendv1.ListOrganizationsResponse{
		Organizations: organizations,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) GetOrganization(ctx context.Context, req *backendv1.GetOrganizationRequest) (*backendv1.GetOrganizationResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	organizationId, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        organizationId,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	return &backendv1.GetOrganizationResponse{Organization: parseOrganization(qProject, qOrg)}, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *backendv1.UpdateOrganizationRequest) (*backendv1.UpdateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	// fetch existing org; this also acts as a permission check
	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	updates := queries.UpdateOrganizationParams{
		ID: orgID,
	}

	updates.DisplayName = qOrg.DisplayName
	if req.Organization.DisplayName != "" {
		updates.DisplayName = req.Organization.DisplayName
	}

	updates.LogInWithGoogle = qOrg.LogInWithGoogle
	if req.Organization.LogInWithGoogle != nil {
		if *req.Organization.LogInWithGoogle && !qProject.LogInWithGoogle {
			return nil, apierror.NewPermissionDeniedError("log in with google is not enabled for this project", fmt.Errorf("log in with google is not enabled for this project"))
		}

		updates.LogInWithGoogle = *req.Organization.LogInWithGoogle
	}

	updates.LogInWithMicrosoft = qOrg.LogInWithMicrosoft
	if req.Organization.LogInWithMicrosoft != nil {
		if *req.Organization.LogInWithMicrosoft && !qProject.LogInWithMicrosoft {
			return nil, apierror.NewPermissionDeniedError("log in with microsoft is not enabled for this project", fmt.Errorf("log in with microsoft is not enabled for this project"))
		}

		updates.LogInWithMicrosoft = *req.Organization.LogInWithMicrosoft
	}

	updates.LogInWithEmail = qOrg.LogInWithEmail
	if req.Organization.LogInWithEmail != nil {
		if *req.Organization.LogInWithEmail && !qProject.LogInWithEmail {
			return nil, apierror.NewPermissionDeniedError("log in with email is not enabled for this project", fmt.Errorf("log in with email is not enabled for this project"))
		}

		updates.LogInWithEmail = *req.Organization.LogInWithEmail
	}

	updates.LogInWithPassword = qOrg.LogInWithPassword
	if req.Organization.LogInWithPassword != nil {
		if *req.Organization.LogInWithPassword && !qProject.LogInWithPassword {
			return nil, apierror.NewPermissionDeniedError("log in with password is not enabled for this project", fmt.Errorf("log in with password is not enabled for this project"))
		}

		updates.LogInWithPassword = *req.Organization.LogInWithPassword
	}

	updates.LogInWithSaml = qOrg.LogInWithSaml
	if req.Organization.LogInWithSaml != nil {
		if *req.Organization.LogInWithSaml && !qProject.LogInWithSaml {
			return nil, apierror.NewPermissionDeniedError("log in with saml is not enabled for this project", fmt.Errorf("log in with saml is not enabled for this project"))
		}

		updates.LogInWithSaml = *req.Organization.LogInWithSaml
	}

	updates.LogInWithAuthenticatorApp = qOrg.LogInWithAuthenticatorApp
	if req.Organization.LogInWithAuthenticatorApp != nil {
		if *req.Organization.LogInWithAuthenticatorApp && !qProject.LogInWithAuthenticatorApp {
			return nil, apierror.NewPermissionDeniedError("log in with authenticator app is not enabled for this project", fmt.Errorf("log in with authenticator app is not enabled for this project"))
		}

		updates.LogInWithAuthenticatorApp = *req.Organization.LogInWithAuthenticatorApp
	}

	updates.LogInWithPasskey = qOrg.LogInWithPasskey
	if req.Organization.LogInWithPasskey != nil {
		if *req.Organization.LogInWithPasskey && !qProject.LogInWithPasskey {
			return nil, apierror.NewPermissionDeniedError("log in with passkey is not enabled for this project", fmt.Errorf("log in with passkey is not enabled for this project"))
		}

		updates.LogInWithPasskey = *req.Organization.LogInWithPasskey
	}

	updates.ScimEnabled = qOrg.ScimEnabled
	if req.Organization.ScimEnabled != nil {
		updates.ScimEnabled = *req.Organization.ScimEnabled
	}

	updates.RequireMfa = qOrg.RequireMfa
	if req.Organization.RequireMfa != nil {
		if *req.Organization.RequireMfa {
			if !updates.LogInWithAuthenticatorApp && !updates.LogInWithPasskey {
				return nil, apierror.NewInvalidArgumentError("require mfa requires log in with authenticator app or passkey to be enabled", fmt.Errorf("require mfa requires log in with authenticator app or passkey to be enabled"))
			}
		}

		updates.RequireMfa = *req.Organization.RequireMfa
	}

	qUpdatedOrg, err := q.UpdateOrganization(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	// Send webhook event
	if err := s.sendSyncOrganizationEvent(ctx, qUpdatedOrg); err != nil {
		return nil, fmt.Errorf("send sync organization event: %w", err)
	}

	return &backendv1.UpdateOrganizationResponse{Organization: parseOrganization(qProject, qUpdatedOrg)}, nil
}

func (s *Store) DeleteOrganization(ctx context.Context, req *backendv1.DeleteOrganizationRequest) (*backendv1.DeleteOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	orgID, err := idformat.Organization.Parse(req.Id)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid organization id", fmt.Errorf("parse organization id: %w", err))
	}

	// authz check
	qOrg, err := q.GetOrganizationByProjectIDAndID(ctx, queries.GetOrganizationByProjectIDAndIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization: %w", err)
	}

	if err := q.DeleteOrganization(ctx, orgID); err != nil {
		return nil, fmt.Errorf("delete organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	// Send webhook event
	if err := s.sendSyncOrganizationEvent(ctx, qOrg); err != nil {
		return nil, fmt.Errorf("send sync organization event: %w", err)
	}

	return &backendv1.DeleteOrganizationResponse{}, nil
}

func (s *Store) DisableOrganizationLogins(ctx context.Context, req *backendv1.DisableOrganizationLoginsRequest) (*backendv1.DisableOrganizationLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.DisableOrganizationLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("lockout organization: %w", err)
	}

	if err := q.RevokeAllOrganizationSessions(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("revoke all organization sessions: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DisableOrganizationLoginsResponse{}, nil
}

func (s *Store) EnableOrganizationLogins(ctx context.Context, req *backendv1.EnableOrganizationLoginsRequest) (*backendv1.EnableOrganizationLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.EnableOrganizationLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("unlock organization: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.EnableOrganizationLoginsResponse{}, nil
}

func (s *Store) sendSyncOrganizationEvent(ctx context.Context, qOrg queries.Organization) error {
	qProjectWebhookSettings, err := s.q.GetProjectWebhookSettings(ctx, authn.ProjectID(ctx))
	if err != nil {
		// We want to ignore this error if the project does not have webhook settings
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("get project by id: %w", err)
	}

	message, err := s.svixClient.Message.Create(ctx, qProjectWebhookSettings.AppID, models.MessageIn{
		EventType: "sync.organization",
		Payload: map[string]interface{}{
			"type": "sync.organization",
			"id":   idformat.Organization.Format(qOrg.ID),
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	slog.InfoContext(ctx, "svix_message_created", "message_id", message.Id, "event_type", message.EventType, "organization_id", idformat.Organization.Format(qOrg.ID))

	return nil
}

func parseOrganization(qProject queries.Project, qOrg queries.Organization) *backendv1.Organization {
	return &backendv1.Organization{
		Id:                        idformat.Organization.Format(qOrg.ID),
		DisplayName:               qOrg.DisplayName,
		CreateTime:                timestamppb.New(*qOrg.CreateTime),
		UpdateTime:                timestamppb.New(*qOrg.UpdateTime),
		LogInWithGoogle:           &qOrg.LogInWithGoogle,
		LogInWithMicrosoft:        &qOrg.LogInWithMicrosoft,
		LogInWithEmail:            &qOrg.LogInWithEmail,
		LogInWithPassword:         &qOrg.LogInWithPassword,
		LogInWithSaml:             &qOrg.LogInWithSaml,
		LogInWithAuthenticatorApp: &qOrg.LogInWithAuthenticatorApp,
		LogInWithPasskey:          &qOrg.LogInWithPasskey,
		RequireMfa:                &qOrg.RequireMfa,
		ScimEnabled:               &qOrg.ScimEnabled,
	}
}
