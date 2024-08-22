package client

import (
	"context"

	"github.com/balobas/auth_service_bln/internal/entity/contract"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type QueryExecer interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	ScanQueryRow(ctx context.Context, dest interface{}, sql string, args ...interface{}) error
	ScanAllQuery(ctx context.Context, dest interface{}, sql string, args ...interface{}) error
}

type Transactor interface {
	BeginTxWithContext(ctx context.Context) (context.Context, contract.Transaction, error)
}

type DB interface {
	QueryExecer
	Transactor
	Ping(ctx context.Context) error
	Close()
}

type ClientDB interface {
	DB() DB
	Close(ctx context.Context) error
}
