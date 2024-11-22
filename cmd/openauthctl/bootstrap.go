package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth-dev/openauth/internal/store"
)

type bootstrapArgs struct {
	Args     args   `cli:"bootstrap,subcmd"`
	Database string `cli:"-d,--database"`
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

	s := store.New(store.NewStoreParams{
		DB: db,
	})

	res, err := s.CreateDogfoodProject(ctx)
	if err != nil {
		return fmt.Errorf("create dogfood project: %w", err)
	}

	fmt.Printf("%s\t%s\t%s\n", res.DogfoodProjectID, res.BootstrapUserEmail, res.BootstrapUserVerySensitivePassword)
	return nil
}
