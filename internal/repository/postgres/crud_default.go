package repositoryPostgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type Row interface {
	IdColumnName() string
	Values() []interface{}
	Columns() []string
	Table() string
	GetId() interface{}
	ScanId(row pgx.Row) error
	Scan(row pgx.Row) error
	ColumnsForUpdate() []string
	ValuesForUpdate() []interface{}
}

func (r *BasePgRepository) Create(ctx context.Context, row Row) error {
	stmt, args, err := sq.Insert(row.Table()).
		PlaceholderFormat(sq.Dollar).
		Columns(row.Columns()...).
		Values(row.Values()...).ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB().Exec(ctx, stmt, args...)
	return err
}

func (r *BasePgRepository) Get(ctx context.Context, row Row) error {
	stmt, args, err := sq.Select(row.Columns()...).
		From(row.Table()).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{row.IdColumnName(): row.GetId()}).ToSql()
	if err != nil {
		return err
	}
	return row.Scan(r.DB().QueryRow(ctx, stmt, args...))
}

func (r *BasePgRepository) Update(ctx context.Context, row Row) error {
	columnsForUpdate := row.ColumnsForUpdate()
	valuesForUpdate := row.ValuesForUpdate()

	sqlBuilder := sq.Update(row.Table()).PlaceholderFormat(sq.Dollar)

	for i := 0; i < len(columnsForUpdate); i++ {
		sqlBuilder = sqlBuilder.Set(columnsForUpdate[i], valuesForUpdate[i])
	}

	sqlBuilder = sqlBuilder.Where(sq.Eq{row.IdColumnName(): row.GetId()})

	stmt, args, err := sqlBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB().Exec(ctx, stmt, args...)
	return err
}

func (r *BasePgRepository) Delete(ctx context.Context, row Row) error {
	stmt, args, err := sq.Delete(row.Table()).PlaceholderFormat(sq.Dollar).Where(sq.Eq{row.IdColumnName(): row.GetId()}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.DB().Exec(ctx, stmt, args...)
	return err
}
