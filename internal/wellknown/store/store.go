package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openauth/openauth/internal/wellknown/store/queries"
)

type Store struct {
	db *pgxpool.Pool
	q  *queries.Queries
}

type NewStoreParams struct {
	DB *pgxpool.Pool
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db: p.DB,
		q:  queries.New(p.DB),
	}

	return store
}
