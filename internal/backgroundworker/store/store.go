package store

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/tesseral-labs/tesseral/internal/backgroundworker/store/queries"
)

type Store struct {
	DB                      *pgxpool.Pool
	Svix                    *svix.Svix
	DirectWebhookHTTPClient *http.Client
}

func (s *Store) q() *queries.Queries {
	return queries.New(s.DB)
}
