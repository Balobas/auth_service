package repo_test

import (
	"context"
	"testing"

	"github.com/balobas/auth_service/internal/client"
	"github.com/balobas/auth_service/internal/client/pg"
)

const pgTestDsn = "host=localhost port=5432 dbname=auth user=auth-user password=auth-password sslmode=disable"

func NewPgClient(t *testing.T, ctx context.Context) client.ClientDB {
	c, err := pg.NewClient(ctx, pgTestDsn)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
