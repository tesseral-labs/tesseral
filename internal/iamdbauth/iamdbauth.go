package iamdbauth

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

type Config struct {
	Region string `conf:"region,noredact"`
	Host   string `conf:"host,noredact"`
	Port   int    `conf:"port,noredact"`
	User   string `conf:"user,noredact"`
	DBName string `conf:"dbname,noredact"`
}

// BuildConnectionString uses AWS IAM database authentication to build a
// connection string to a given database.
func BuildConnectionString(ctx context.Context, c Config) (string, error) {
	awscfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("load default aws config: %w", err)
	}

	endpoint := fmt.Sprintf("%s:%d", c.Host, c.Port)
	authToken, err := auth.BuildAuthToken(ctx, endpoint, c.Region, c.User, awscfg.Credentials)
	if err != nil {
		return "", fmt.Errorf("build rds auth token: %w", err)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.Host, c.Port, c.User, authToken, c.DBName), nil
}
