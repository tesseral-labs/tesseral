package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/kms"
	"github.com/tesseral-labs/tesseral/internal/store"
)

type bootstrapArgs struct {
	Args          args   `cli:"bootstrap,subcmd"`
	Database      string `cli:"--database"`
	ConsoleDomain string `cli:"--console-domain"`
	VaultDomain   string `cli:"--vault-domain"`
	RootUserEmail string `cli:"--root-user-email"`

	SessionSigningKeysKMSBackend                 string `cli:"--session-signing-keys-kms-backend"`
	SessionSigningKeysKMSAWSKMSV1KeyID           string `cli:"--session-signing-keys-kms-aws-kms-v1-key-id"`
	SessionSigningKeysKMSAWSKMSV1KMSBaseEndpoint string `cli:"--session-signing-keys-kms-aws-kms-v1-kms-base-endpoint"`
	SessionSigningKeysKMSGCPKMSV1KeyName         string `cli:"--session-signing-keys-kms-gcp-kms-v1-key-name"`
}

func (bootstrapArgs) Description() string {
	return "Bootstrap a Tesseral database"
}

func (bootstrapArgs) ExtendedDescription() string {
	return strings.TrimSpace(`
Bootstrap a Tesseral database.

Outputs, tab-separated, a project ID, an email, and a very sensitive password.

The project ID is the console project ID. The email and password are a login
method for a "root" user, an owner of the console project.

Rotate or delete the root password before deploying this Tesseral instance in
production.
`)
}

func bootstrap(ctx context.Context, args bootstrapArgs) error {
	db, err := pgxpool.New(context.Background(), args.Database)
	if err != nil {
		panic(err)
	}

	sessionSigningKeysKMS, err := kms.New(ctx, kms.Config{
		Backend:                 args.SessionSigningKeysKMSBackend,
		AWSKMSV1KeyID:           args.SessionSigningKeysKMSAWSKMSV1KeyID,
		AWSKMSV1KMSBaseEndpoint: args.SessionSigningKeysKMSAWSKMSV1KMSBaseEndpoint,
		GCPKMSV1KeyName:         args.SessionSigningKeysKMSGCPKMSV1KeyName,
	})
	if err != nil {
		panic(fmt.Errorf("create session signing keys kms: %w", err))
	}

	s := store.New(store.NewStoreParams{
		DB:                   db,
		SessionSigningKeyKMS: sessionSigningKeysKMS,
	})

	res, err := s.CreateConsoleProject(ctx, &store.CreateConsoleProjectRequest{
		RootUserEmail: args.RootUserEmail,
		ConsoleDomain: args.ConsoleDomain,
		VaultDomain:   args.VaultDomain,
	})
	if err != nil {
		return fmt.Errorf("create console project: %w", err)
	}

	fmt.Printf(
		"%s\t%s\t%s\n",
		res.ConsoleProjectID,
		res.BootstrapUserEmail,
		res.BootstrapUserVerySensitivePassword,
	)
	return nil
}
