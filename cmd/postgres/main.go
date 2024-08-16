package main

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const DbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"

func main() {

	ctx := context.Background()

	conn, err := pgxpool.Connect(ctx, DbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	sqlBuilder := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "token").
		Values(15, "AAAAAAAAA")

	stmt, args, err := sqlBuilder.ToSql()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(stmt)

	// ct, err := conn.Exec(ctx, stmt, args...)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(ct.RowsAffected())

	b := sq.Select("id", "token").From("auth").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": 12})

	// b := sq.Select("a.id", "a.token").
	// 	From("auth a").
	// 	LeftJoin("hui h ON a.id = h.id").
	// 	PlaceholderFormat(sq.Dollar).
	// 	Where("h.pizda = ?", 12)

	stmt, args, err = b.ToSql()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(stmt)

	row := conn.QueryRow(ctx, stmt, args...)
	var id int64
	var token string
	if err := row.Scan(&id, &token); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(id, token)

}
