package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth/openauth/internal/loadenv"
	"github.com/openauth/openauth/internal/store"
)

type bootstrapArgs struct {
	Args                               args   `cli:"bootstrap,subcmd"`
	Database                           string `cli:"-d,--database"`
	KMSEndpoint                        string `cli:"-k,--kms-endpoint"`
	IntermediateSessionSigningKMSKeyID string `cli:"-i,--intermediate-session-kms-key-id"`
	SessionSigningKMSKeyID             string `cli:"-s,--session-kms-key-id"`
}

func (_ bootstrapArgs) Description() string {
	return "Bootstrap an OpenAuth database"
}

func (_ bootstrapArgs) ExtendedDescription() string {
	return strings.TrimSpace(`
Bootstrap an OpenAuth database.

Outputs, tab-separated, a project ID, an email, and a very sensitive password.

The project ID is the bootstrap ("dogfood") project ID. The email and password
are a login method for an admin user in that project.

Delete this admin user before deploying this OpenAuth instance in production.
`)
}

func bootstrap(ctx context.Context, args bootstrapArgs) error {
	db, err := pgxpool.New(context.Background(), args.Database)
	if err != nil {
		panic(err)
	}

	loadenv.LoadEnv()

	awsConf, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("load aws config: %w", err))
	}

	kms_ := kms.NewFromConfig(awsConf, func(o *kms.Options) {
		if args.KMSEndpoint != "" {
			o.BaseEndpoint = &args.KMSEndpoint
		}
	})

	s := store.New(store.NewStoreParams{
		DB:                                    db,
		IntermediateSessionSigningKeyKMSKeyID: args.IntermediateSessionSigningKMSKeyID,
		KMS:                                   kms_,
		SessionSigningKeyKmsKeyID:             args.SessionSigningKMSKeyID,
	})

	res, err := s.CreateDogfoodProject(ctx)
	if err != nil {
		return fmt.Errorf("create dogfood project: %w", err)
	}

	fmt.Printf(
		"%s\t%s\t%s\n",
		res.DogfoodProjectID,
		res.BootstrapUserEmail,
		res.BootstrapUserVerySensitivePassword,
	)
	return nil
}
