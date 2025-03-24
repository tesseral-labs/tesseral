package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/store"
)

type bootstrapArgs struct {
	Args                       args   `cli:"bootstrap,subcmd"`
	Database                   string `cli:"--database"`
	KMSEndpoint                string `cli:"--kms-endpoint"`
	SessionSigningKMSKeyID     string `cli:"--session-kms-key-id"`
	AuthAppsRootDomain         string `cli:"--auth-apps-root-domain"`
	RootUserEmail              string `cli:"--root-user-email"`
	DogfoodProjectRedirectURI  string `cli:"--dogfood-project-redirect-uri"`
	DogfoodProjectCookieDomain string `cli:"--dogfood-project-cookie-domain"`
}

func (_ bootstrapArgs) Description() string {
	return "Bootstrap a Tesseral database"
}

func (_ bootstrapArgs) ExtendedDescription() string {
	return strings.TrimSpace(`
Bootstrap a Tesseral database.

Outputs, tab-separated, a project ID, an email, and a very sensitive password.

The project ID is the bootstrap ("dogfood") project ID. The email and password
are a login method for an admin user in that project.

Delete this admin user before deploying this Tesseral instance in production.
`)
}

func bootstrap(ctx context.Context, args bootstrapArgs) error {
	db, err := pgxpool.New(context.Background(), args.Database)
	if err != nil {
		panic(err)
	}

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
		DB:                        db,
		KMS:                       kms_,
		SessionSigningKeyKmsKeyID: args.SessionSigningKMSKeyID,
	})

	res, err := s.CreateDogfoodProject(ctx, &store.CreateDogfoodProjectRequest{
		AuthAppsRootDomain: args.AuthAppsRootDomain,
		RootUserEmail:      args.RootUserEmail,
		RedirectURI:        args.DogfoodProjectRedirectURI,
		CookieDomain:       args.DogfoodProjectCookieDomain,
	})
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
