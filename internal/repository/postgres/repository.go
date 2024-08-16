package repositoryPostgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

type RepositoryConfig interface {
	DSN() string
}

func New(ctx context.Context, config RepositoryConfig) (*Repository, error) {
	pool, err := pgxpool.Connect(ctx, config.DSN())
	if err != nil {
		return nil, err
	}

	return &Repository{pool: pool}, nil
}

func (r *Repository) Close() {
	r.pool.Close()
}
