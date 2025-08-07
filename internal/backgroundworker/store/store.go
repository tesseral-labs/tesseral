package store

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/jackc/pgx/v5/pgxpool"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store/queries"
)

type Store struct {
	DB                      *pgxpool.Pool
	Svix                    *svix.Svix
	DirectWebhookHTTPClient *http.Client
	SES                     *sesv2.Client
	ConsoleProjectID        string
	ConsoleDomain           string
}

func (s *Store) q() *queries.Queries {
	return queries.New(s.DB)
}
