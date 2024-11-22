package main

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/ucarion/cli"
)

func main() {
	cli.Run(context.Background(), version, force, up, bootstrap)
}

type args struct {
}

func (_ args) ExtendedDescription() string {
	return "Control the openauth database"
}
