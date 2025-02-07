package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetProject(ctx context.Context, req *backendv1.GetProjectRequest) (*backendv1.GetProjectResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	project, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, err
	}

	qProjectPasskeyRPIDs, err := q.GetProjectPasskeyRPIDs(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project passkey rp ids: %w", err)
	}

	return &backendv1.GetProjectResponse{Project: parseProject(&project, qProjectPasskeyRPIDs)}, nil
}

func (s *Store) DisableProjectLogins(ctx context.Context, req *backendv1.DisableProjectLoginsRequest) (*backendv1.DisableProjectLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.DisableProjectLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("lockout project: %w", err)
	}

	if err := q.RevokeAllProjectSessions(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("revoke all project sessions: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DisableProjectLoginsResponse{}, nil
}

func (s *Store) EnableProjectLogins(ctx context.Context, req *backendv1.EnableProjectLoginsRequest) (*backendv1.EnableProjectLoginsResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := q.EnableProjectLogins(ctx, authn.ProjectID(ctx)); err != nil {
		return nil, fmt.Errorf("unlock project: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.EnableProjectLoginsResponse{}, nil
}

func (s *Store) UpdateProject(ctx context.Context, req *backendv1.UpdateProjectRequest) (*backendv1.UpdateProjectResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	// fetch project outside a transaction, so that we can carry out KMS
	// operations; we can live with possibility of conflicting concurrent writes
	qProject, err := s.q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %w", err))
		}

		return nil, fmt.Errorf("get project by id: %w", err)
	}

	updates := queries.UpdateProjectParams{
		ID: qProject.ID,
	}

	updates.DisplayName = qProject.DisplayName
	if req.Project.DisplayName != "" {
		updates.DisplayName = req.Project.DisplayName
	}

	updates.GoogleOauthClientID = qProject.GoogleOauthClientID
	if req.Project.GoogleOauthClientId != "" {
		updates.GoogleOauthClientID = &req.Project.GoogleOauthClientId
	}

	updates.GoogleOauthClientSecretCiphertext = qProject.GoogleOauthClientSecretCiphertext
	if req.Project.GoogleOauthClientSecret != "" {
		encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               &s.googleOAuthClientSecretsKMSKeyID,
			Plaintext:           []byte(req.Project.GoogleOauthClientSecret),
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		})
		if err != nil {
			return nil, fmt.Errorf("encrypt google oauth client secret: %w", err)
		}

		updates.GoogleOauthClientSecretCiphertext = encryptRes.CiphertextBlob
	}

	updates.MicrosoftOauthClientID = qProject.MicrosoftOauthClientID
	if req.Project.MicrosoftOauthClientId != "" {
		updates.MicrosoftOauthClientID = &req.Project.MicrosoftOauthClientId
	}

	updates.MicrosoftOauthClientSecretCiphertext = qProject.MicrosoftOauthClientSecretCiphertext
	if req.Project.MicrosoftOauthClientSecret != "" {
		encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
			KeyId:               &s.microsoftOAuthClientSecretsKMSKeyID,
			Plaintext:           []byte(req.Project.MicrosoftOauthClientSecret),
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		})
		if err != nil {
			return nil, fmt.Errorf("encrypt microsoft oauth client secret: %w", err)
		}

		updates.MicrosoftOauthClientSecretCiphertext = encryptRes.CiphertextBlob
	}

	updates.LogInWithGoogle = qProject.LogInWithGoogle
	if req.Project.LogInWithGoogle != nil {
		// todo: validate that google is configured?
		updates.LogInWithGoogle = *req.Project.LogInWithGoogle
	}

	updates.LogInWithMicrosoft = qProject.LogInWithMicrosoft
	if req.Project.LogInWithMicrosoft != nil {
		// todo: validate that microsoft is configured?
		updates.LogInWithMicrosoft = *req.Project.LogInWithMicrosoft
	}

	updates.LogInWithEmail = qProject.LogInWithEmail
	if req.Project.LogInWithEmail != nil {
		updates.LogInWithEmail = *req.Project.LogInWithEmail
	}

	updates.LogInWithPassword = qProject.LogInWithPassword
	if req.Project.LogInWithPassword != nil {
		updates.LogInWithPassword = *req.Project.LogInWithPassword
	}

	updates.LogInWithSaml = qProject.LogInWithSaml
	if req.Project.LogInWithSaml != nil {
		updates.LogInWithSaml = *req.Project.LogInWithSaml
	}

	updates.LogInWithAuthenticatorApp = qProject.LogInWithAuthenticatorApp
	if req.Project.LogInWithAuthenticatorApp != nil {
		updates.LogInWithAuthenticatorApp = *req.Project.LogInWithAuthenticatorApp
	}

	updates.LogInWithPasskey = qProject.LogInWithPasskey
	if req.Project.LogInWithPasskey != nil {
		updates.LogInWithPasskey = *req.Project.LogInWithPasskey
	}

	updates.CustomAuthDomain = qProject.CustomAuthDomain
	// TODO enable when we have a use for custom domains from app
	//if req.Project.CustomDomain != nil {
	//	updates.CustomDomain = req.Project.CustomDomain
	//}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qUpdatedProject, err := q.UpdateProject(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	if !qUpdatedProject.LogInWithGoogle {
		if _, err := q.DisableProjectOrganizationsLogInWithGoogle(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with Google: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithMicrosoft {
		if _, err := q.DisableProjectOrganizationsLogInWithMicrosoft(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with Microsoft: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithEmail {
		if _, err := q.DisableProjectOrganizationsLogInWithEmail(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with email: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithPassword {
		if _, err := q.DisableProjectOrganizationsLogInWithPassword(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with password: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithSaml {
		if _, err := q.DisableProjectOrganizationsLogInWithSAML(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with SAML: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithAuthenticatorApp {
		if _, err := q.DisableProjectOrganizationsLogInWithAuthenticatorApp(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with authenticator app: %w", err)
		}
	}

	if !qUpdatedProject.LogInWithPasskey {
		if _, err := q.DisableProjectOrganizationsLogInWithPasskey(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable project organizations log in with passkey: %w", err)
		}
	}

	// only update project passkey RP IDs if mentioned in request
	if len(req.Project.PasskeyRpIds) > 0 {
		// dedupe RP IDs, and always include the current vault domain in list
		passkeyRPIDs := map[string]struct{}{*qUpdatedProject.AuthDomain: {}}
		for _, rpID := range req.Project.PasskeyRpIds {
			passkeyRPIDs[rpID] = struct{}{}
		}

		if err := q.DeleteProjectPasskeyRPIDs(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("delete project passkey rp ids: %w", err)
		}

		for rpID := range passkeyRPIDs {
			if _, err := q.CreateProjectPasskeyRPID(ctx, queries.CreateProjectPasskeyRPIDParams{
				ProjectID: authn.ProjectID(ctx),
				RpID:      rpID,
			}); err != nil {
				return nil, fmt.Errorf("create project passkey rp id: %w", err)
			}
		}

		// disable all existing passkeys in project outside new list of RP IDs
		if err := q.DisablePasskeysOutsideProjectRPIDs(ctx, authn.ProjectID(ctx)); err != nil {
			return nil, fmt.Errorf("disable passkeys outside project rp ids: %w", err)
		}
	}

	qProjectPasskeyRPIDs, err := q.GetProjectPasskeyRPIDs(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project passkey rp ids: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateProjectResponse{Project: parseProject(&qUpdatedProject, qProjectPasskeyRPIDs)}, nil
}

func parseProject(qProject *queries.Project, qProjectPasskeyRPIDs []queries.ProjectPasskeyRpID) *backendv1.Project {
	// sanity check
	for _, qProjectPasskeyRPID := range qProjectPasskeyRPIDs {
		if qProjectPasskeyRPID.ProjectID != qProject.ID {
			panic(fmt.Errorf("project passkey rp id project id mismatch: %s != %s", qProjectPasskeyRPID.ProjectID, qProject.ID))
		}
	}

	authDomain := derefOrEmpty(qProject.AuthDomain)
	if qProject.CustomAuthDomain != nil {
		authDomain = *qProject.CustomAuthDomain
	}

	var passkeyRPIDs []string
	for _, qProjectPasskeyRPID := range qProjectPasskeyRPIDs {
		passkeyRPIDs = append(passkeyRPIDs, qProjectPasskeyRPID.RpID)
	}

	return &backendv1.Project{
		Id:                        idformat.Project.Format(qProject.ID),
		DisplayName:               qProject.DisplayName,
		CreateTime:                timestamppb.New(*qProject.CreateTime),
		UpdateTime:                timestamppb.New(*qProject.UpdateTime),
		LogInWithGoogle:           &qProject.LogInWithMicrosoft,
		LogInWithMicrosoft:        &qProject.LogInWithMicrosoft,
		LogInWithEmail:            &qProject.LogInWithEmail,
		LogInWithPassword:         &qProject.LogInWithPassword,
		LogInWithSaml:             &qProject.LogInWithSaml,
		LogInWithAuthenticatorApp: &qProject.LogInWithAuthenticatorApp,
		LogInWithPasskey:          &qProject.LogInWithPasskey,
		GoogleOauthClientId:       derefOrEmpty(qProject.GoogleOauthClientID),
		MicrosoftOauthClientId:    derefOrEmpty(qProject.MicrosoftOauthClientID),
		AuthDomain:                &authDomain,
		PasskeyRpIds:              passkeyRPIDs,
	}
}
