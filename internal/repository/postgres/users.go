package repositoryPostgres

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
)

type CreateUserParams interface {
	GetName() string
	GetEmail() string
	GetPassword() string
	GetPasswordConfirm() string
}

type UserParams interface {
	GetId() int64
	CreateUserParams
}

type Role interface {
	String() string
}

func (r *Repository) CreateUser(ctx context.Context, user CreateUserParams, role Role) (int64, error) {
	sqlBuilder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "role", "created_at", "updated_at").
		Values(
			user.GetName(),
			user.GetEmail(),
			user.GetPassword(),
			user.GetPasswordConfirm(),
			role.String(),
			pgtype.Timestamp{
				Time:   time.Now().UTC(),
				Status: pgtype.Present,
			},
			pgtype.Timestamp{
				Time:   time.Now().UTC(),
				Status: pgtype.Present,
			},
		).Suffix("RETURNING id")

	stmt, args, err := sqlBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	row := r.pool.QueryRow(ctx, stmt, args...)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	log.Printf("successfuly create user id: %d, name: %s, email: %s \n", id, user.GetName(), user.GetEmail())

	return id, nil
}

type RepositoryUser struct {
	Id              int64
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
	Role            string
	CreatedAt       int64
	UpdatedAt       int64
}

func (r *Repository) GetUser(ctx context.Context, id int64) (*RepositoryUser, error) {
	sqlBuilder := sq.Select("id", "name", "email", "password", "password_confirm", "role", "created_at", "updated_at").
		From("users").PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})
	stmt, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var createdAt, updatedAt pgtype.Timestamp

	var user RepositoryUser
	if err := r.pool.QueryRow(ctx, stmt, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.ConfirmPassword,
		&user.Role,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}

	user.CreatedAt = createdAt.Time.Unix()
	user.UpdatedAt = updatedAt.Time.Unix()

	return &user, nil
}

type UpdateUserParams interface {
	GetName() string
	GetEmail() string
	GetId() int64
}

func (r *Repository) UpdateUser(ctx context.Context, id int64, name string, email string) error {
	sqlBuilder := sq.Update("users").PlaceholderFormat(sq.Dollar).
		Set("name", name).
		Set("email", email).
		Set("updated_at", pgtype.Timestamp{
			Time:   time.Now().UTC(),
			Status: pgtype.Present,
		}).
		Where(sq.Eq{"id": id})

	stmt, args, err := sqlBuilder.ToSql()
	if err != nil {
		return err
	}

	tag, err := r.pool.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	log.Printf("update user with uid: %d affected rows: %d\n", id, tag.RowsAffected())

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	stmt, args, err := sq.Delete("users").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	fmt.Println(stmt)

	tag, err := r.pool.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	log.Printf("delete user with uid: %d affected rows: %d\n", id, tag.RowsAffected())

	return nil
}
