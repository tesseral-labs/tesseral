package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/svix/svix-webhooks/go/models"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetOrganization(ctx context.Context, req *frontendv1.GetOrganizationRequest) (*frontendv1.GetOrganizationResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrganization, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
	}

	return &frontendv1.GetOrganizationResponse{Organization: parseOrganization(qProject, qOrganization)}, nil
}

func (s *Store) UpdateOrganization(ctx context.Context, req *frontendv1.UpdateOrganizationRequest) (*frontendv1.UpdateOrganizationResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := s.validateIsOwner(ctx); err != nil {
		return nil, fmt.Errorf("validate is owner: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetOrganizationByID(ctx, authn.OrganizationID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("organization not found", fmt.Errorf("get organization by id: %w", err))
		}

		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	updates := queries.UpdateOrganizationParams{
		ID: authn.OrganizationID(ctx),
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

	updates.LogInWithGithub = qOrg.LogInWithGithub
	if req.Organization.LogInWithGithub != nil {
		if *req.Organization.LogInWithGithub && !qProject.LogInWithGithub {
			return nil, apierror.NewPermissionDeniedError("log in with github is not enabled for this project", fmt.Errorf("log in with github is not enabled for this project"))
		}

		updates.LogInWithGithub = *req.Organization.LogInWithGithub
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
		return nil, fmt.Errorf("update organization: %w", fmt.Errorf("update organization: %w", err))
	}

	// Commit the transaction
	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	organization := parseOrganization(qProject, qUpdatedOrg)
	previousOrganization := parseOrganization(qProject, qOrg)
	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.organizations.update",
		EventDetails: map[string]any{
			"organization":         organization,
			"previousOrganization": previousOrganization,
		},
		ResourceType: queries.AuditLogEventResourceTypeOrganization,
		ResourceID:   &qOrg.ID,
	}); err != nil {
		return nil, fmt.Errorf("create audit log event: %w", err)
	}

	// send sync organization event
	if err := s.sendSyncOrganizationEvent(ctx, qUpdatedOrg); err != nil {
		return nil, fmt.Errorf("send sync organization event: %w", err)
	}

	return &frontendv1.UpdateOrganizationResponse{
		Organization: organization,
	}, nil
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
			"type":           "sync.organization",
			"organizationId": idformat.Organization.Format(qOrg.ID),
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	slog.InfoContext(ctx, "svix_message_created", "message_id", message.Id, "event_type", message.EventType, "organization_id", idformat.Organization.Format(qOrg.ID))

	return nil
}

func parseOrganization(qProject queries.Project, qOrg queries.Organization) *frontendv1.Organization {
	return &frontendv1.Organization{
		Id:                        idformat.Organization.Format(qOrg.ID),
		DisplayName:               qOrg.DisplayName,
		CreateTime:                timestamppb.New(*qOrg.CreateTime),
		UpdateTime:                timestamppb.New(*qOrg.UpdateTime),
		LogInWithGoogle:           &qOrg.LogInWithGoogle,
		LogInWithMicrosoft:        &qOrg.LogInWithMicrosoft,
		LogInWithGithub:           &qOrg.LogInWithGithub,
		LogInWithEmail:            &qOrg.LogInWithEmail,
		LogInWithPassword:         &qOrg.LogInWithPassword,
		LogInWithSaml:             &qOrg.LogInWithSaml,
		LogInWithAuthenticatorApp: &qOrg.LogInWithAuthenticatorApp,
		LogInWithPasskey:          &qOrg.LogInWithPasskey,
		RequireMfa:                &qOrg.RequireMfa,
		GoogleHostedDomains:       nil, // TODO
		MicrosoftTenantIds:        nil, // TODO,
		CustomRolesEnabled:        qOrg.CustomRolesEnabled,
		ApiKeysEnabled:            qOrg.ApiKeysEnabled && qProject.ApiKeysEnabled && qProject.EntitledBackendApiKeys,
		ScimEnabled:               qOrg.ScimEnabled,
	}
}

// validateIsOwner returns an error if the current user is not an owner of the
// organization.
func (s *Store) validateIsOwner(ctx context.Context) error {
	qUser, err := s.q.GetUserByID(ctx, authn.UserID(ctx))
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	if !qUser.IsOwner {
		return apierror.NewPermissionDeniedError("user must be an owner of the organization", fmt.Errorf("user is not an owner"))
	}
	return nil
}
