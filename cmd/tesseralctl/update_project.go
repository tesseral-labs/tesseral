package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/store"
)

type updateProjectArgs struct {
	Args            args   `cli:"update-project,subcmd"`
	ProjectID       string `cli:"project-id"`
	Database        string `cli:"--database"`
	VaultDomain     string `cli:"--vault-domain"`
	LogoURL         string `cli:"--logo-url"`
	DarkModeLogoURL string `cli:"--dark-mode-logo-url"`
}

func (updateProjectArgs) Description() string {
	return "Update a Tesseral Project directly"
}

func (updateProjectArgs) ExtendedDescription() string {
	return strings.TrimSpace(`
Update a Tesseral Project directly, without the Tesseral Console or Backend API.
`)
}

func updateProject(ctx context.Context, args updateProjectArgs) error {
	db, err := pgxpool.New(context.Background(), args.Database)
	if err != nil {
		panic(err)
	}

	s := store.New(store.NewStoreParams{
		DB: db,
	})

	if err := s.UpdateProject(ctx, &store.UpdateProjectRequest{
		ProjectID:       args.ProjectID,
		VaultDomain:     args.VaultDomain,
		LogoURL:         args.LogoURL,
		DarkModeLogoURL: args.DarkModeLogoURL,
	}); err != nil {
		return fmt.Errorf("store: %w", err)
	}

	return nil
}
