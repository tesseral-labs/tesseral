package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/kms"
	"github.com/tesseral-labs/tesseral/internal/store"
)

type createProjectArgs struct {
	Args             args   `cli:"create-project,subcmd"`
	Database         string `cli:"--database"`
	ConsoleProjectID string `cli:"--console-project-id"`
	ProjectName      string `cli:"--project-name"`
	OwnerUserEmail   string `cli:"--owner-user-email"`
	VaultDomain      string `cli:"--vault-domain"`
	LogoURL          string `cli:"--logo-url"`
	DarkModeLogoURL  string `cli:"--dark-mode-logo-url"`

	SessionSigningKeysKMSBackend                 string `cli:"--session-signing-keys-kms-backend"`
	SessionSigningKeysKMSAWSKMSV1KeyID           string `cli:"--session-signing-keys-kms-aws-kms-v1-key-id"`
	SessionSigningKeysKMSAWSKMSV1KMSBaseEndpoint string `cli:"--session-signing-keys-kms-aws-kms-v1-kms-base-endpoint"`
	SessionSigningKeysKMSGCPKMSV1KeyName         string `cli:"--session-signing-keys-kms-gcp-kms-v1-key-name"`
}

func (createProjectArgs) Description() string {
	return "Create a Tesseral Project directly"
}

func (createProjectArgs) ExtendedDescription() string {
	return strings.TrimSpace(`
Create a Tesseral Project directly, without the Tesseral Console or Backend API.
`)
}

func createProject(ctx context.Context, args createProjectArgs) error {
	db, err := pgxpool.New(context.Background(), args.Database)
	if err != nil {
		return fmt.Errorf("create db pool: %w", err)
	}

	sessionSigningKeysKMS, err := kms.New(ctx, kms.Config{
		Backend:                 args.SessionSigningKeysKMSBackend,
		AWSKMSV1KeyID:           args.SessionSigningKeysKMSAWSKMSV1KeyID,
		AWSKMSV1KMSBaseEndpoint: args.SessionSigningKeysKMSAWSKMSV1KMSBaseEndpoint,
		GCPKMSV1KeyName:         args.SessionSigningKeysKMSGCPKMSV1KeyName,
	})
	if err != nil {
		return fmt.Errorf("create session signing keys kms: %w", err)
	}

	s := store.New(store.NewStoreParams{
		DB:                   db,
		SessionSigningKeyKMS: sessionSigningKeysKMS,
	})

	res, err := s.CreateProject(ctx, &store.CreateProjectRequest{
		ConsoleProjectID: args.ConsoleProjectID,
		ProjectName:      args.ProjectName,
		OwnerUserEmail:   args.OwnerUserEmail,
		VaultDomain:      args.VaultDomain,
		LogoURL:          args.LogoURL,
		DarkModeLogoURL:  args.DarkModeLogoURL,
	})
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	fmt.Printf(
		"%s\t%s\n",
		res.ProjectID,
		res.OwnerPassword,
	)

	return nil
}
