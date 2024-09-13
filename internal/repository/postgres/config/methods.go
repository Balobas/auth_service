package config

import (
	"context"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *ConfigRepository) UpdateConfig(ctx context.Context, cfg map[string]json.RawMessage) error {

	sql := sq.Insert(pgEntity.ConfigTableName).Columns("key", "value").PlaceholderFormat(sq.Dollar)

	for k, v := range cfg {
		sql = sql.Values(k, v)
	}
	stmt, args, err := sql.ToSql()
	if err != nil {
		return errors.WithStack(err)
	}

	stmt += " ON CONFLICT DO UPDATE SET value=EXCLUDED.value"

	if _, err = r.client.DB().Exec(ctx, stmt, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ConfigRepository) GetConfig(ctx context.Context) (map[string]json.RawMessage, error) {
	stmt, args, err := sq.Select("key", "value").
		From(pgEntity.ConfigTableName).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := r.client.DB().Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make(map[string]json.RawMessage, rows.CommandTag().RowsAffected())

	for rows.Next() {
		var key string
		var value json.RawMessage

		if err := rows.Scan(&key, &value); err != nil {
			return nil, errors.WithStack(err)
		}

		res[key] = value
	}

	return res, nil
}
