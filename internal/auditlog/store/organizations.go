package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/auditlog/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetOrganization(ctx context.Context, db queries.DBTX, id uuid.UUID) (*auditlogv1.Organization, error) {
	qOrganization, err := queries.New(db).GetOrganization(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	return &auditlogv1.Organization{
		Id:                        idformat.Organization.Format(qOrganization.ID),
		DisplayName:               qOrganization.DisplayName,
		CreateTime:                timestamppb.New(*qOrganization.CreateTime),
		UpdateTime:                timestamppb.New(*qOrganization.UpdateTime),
		LogInWithPassword:         &qOrganization.LogInWithPassword,
		LogInWithGoogle:           &qOrganization.LogInWithGoogle,
		LogInWithMicrosoft:        &qOrganization.LogInWithMicrosoft,
		LogInWithSaml:             &qOrganization.LogInWithSaml,
		ScimEnabled:               &qOrganization.ScimEnabled,
		LogInWithAuthenticatorApp: &qOrganization.LogInWithAuthenticatorApp,
		LogInWithPasskey:          &qOrganization.LogInWithPasskey,
		RequireMfa:                &qOrganization.RequireMfa,
		LogInWithEmail:            &qOrganization.LogInWithEmail,
		CustomRolesEnabled:        &qOrganization.CustomRolesEnabled,
		ApiKeysEnabled:            &qOrganization.ApiKeysEnabled,
		LogInWithGithub:           &qOrganization.LogInWithGithub,
	}, nil
}
