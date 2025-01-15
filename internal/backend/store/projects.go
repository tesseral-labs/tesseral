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
	"github.com/openauth/openauth/internal/shared/apierror"
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

	return &backendv1.GetProjectResponse{Project: parseProject(&project)}, nil
}

func (s *Store) UpdateProject(ctx context.Context, req *backendv1.UpdateProjectRequest) (*backendv1.UpdateProjectResponse, error) {
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

	updates.LogInWithGoogleEnabled = qProject.LogInWithGoogleEnabled
	if req.Project.LogInWithGoogleEnabled != nil {
		// todo: validate that google is configured?
		updates.LogInWithGoogleEnabled = *req.Project.LogInWithGoogleEnabled
	}

	updates.LogInWithMicrosoftEnabled = qProject.LogInWithMicrosoftEnabled
	if req.Project.LogInWithMicrosoftEnabled != nil {
		// todo: validate that microsoft is configured?
		updates.LogInWithMicrosoftEnabled = *req.Project.LogInWithMicrosoftEnabled
	}

	updates.LogInWithPasswordEnabled = qProject.LogInWithPasswordEnabled
	if req.Project.LogInWithPasswordEnabled != nil {
		updates.LogInWithPasswordEnabled = *req.Project.LogInWithPasswordEnabled
	}

	updates.OrganizationsSamlEnabledDefault = qProject.OrganizationsSamlEnabledDefault
	if req.Project.OrganizationsSamlEnabledDefault != nil {
		updates.OrganizationsSamlEnabledDefault = *req.Project.OrganizationsSamlEnabledDefault
	}

	updates.OrganizationsScimEnabledDefault = qProject.OrganizationsScimEnabledDefault
	if req.Project.OrganizationsScimEnabledDefault != nil {
		updates.OrganizationsScimEnabledDefault = *req.Project.OrganizationsScimEnabledDefault
	}

	// TODO: Re-enable this when we have a need for custom domains on the app side of things

	// updates.CustomDomain = qProject.CustomDomain
	// if req.Project.CustomDomain != nil {
	// 	updates.CustomDomain = req.Project.CustomDomain
	// }

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qUpdatedProject, err := q.UpdateProject(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateProjectResponse{Project: parseProject(&qUpdatedProject)}, nil
}

func parseProject(qProject *queries.Project) *backendv1.Project {
	authDomain := derefOrEmpty(qProject.AuthDomain)
	if qProject.CustomAuthDomain != nil {
		authDomain = *qProject.CustomAuthDomain
	}

	return &backendv1.Project{
		Id:                              idformat.Project.Format(qProject.ID),
		DisplayName:                     qProject.DisplayName,
		CreateTime:                      timestamppb.New(*qProject.CreateTime),
		UpdateTime:                      timestamppb.New(*qProject.UpdateTime),
		LogInWithPasswordEnabled:        &qProject.LogInWithPasswordEnabled,
		LogInWithGoogleEnabled:          &qProject.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled:       &qProject.LogInWithMicrosoftEnabled,
		GoogleOauthClientId:             derefOrEmpty(qProject.GoogleOauthClientID),
		MicrosoftOauthClientId:          derefOrEmpty(qProject.MicrosoftOauthClientID),
		OrganizationsSamlEnabledDefault: &qProject.OrganizationsSamlEnabledDefault,
		OrganizationsScimEnabledDefault: &qProject.OrganizationsScimEnabledDefault,
		AuthDomain:                      &authDomain,
	}
}
