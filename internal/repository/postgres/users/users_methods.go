package repositoryUsers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user entity.User) error {
	userRow := pgEntity.NewUserRow().FromEntity(user)
	if err := r.Create(ctx, userRow); err != nil {
		return errors.Wrapf(err, "failed to create user %s", user.Uid)
	}

	log.Printf("successfuly create user uid: %s, email: %s", user.Uid, user.Email)
	return nil
}

func (r *UsersRepository) GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, bool, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Uid: uid})

	if err := r.GetOne(ctx, userRow, userRow.ConditionUserUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, false, nil
		}
		log.Printf("failed to get user by uid %v", err)
		return entity.User{}, false, errors.WithStack(err)
	}

	log.Printf("successfully get user by uid %s", uid)
	return userRow.ToEntity(), true, nil
}

func (r *UsersRepository) GetAdminUsers(ctx context.Context) ([]entity.User, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{})
	userRows := pgEntity.NewUserRows()

	stmt := fmt.Sprintf(
		`select %s from users inner join user_roles on users.uid = user_roles.user_uid 
		where user_roles.role = $1`, strings.Join(userRow.Columns(), ", "))

	rows, err := r.DB().Query(ctx, stmt, entity.UserRoleAdmin)
	if err != nil {
		log.Printf("failed to get admin users %v", err)
		return []entity.User{}, errors.WithStack(err)
	}
	defer rows.Close()

	if err := userRows.ScanAll(rows); err != nil {
		log.Printf("failed to scan admin users %v", err)
		return []entity.User{}, errors.WithStack(err)
	}

	log.Printf("successfully get admin users")
	return userRows.ToEntity(), nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (entity.User, bool, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Email: email})

	if err := r.GetOne(ctx, userRow, userRow.ConditionEmailEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, false, nil
		}
		log.Printf("failed to get user by uid %v", err)
		return entity.User{}, false, errors.WithStack(err)
	}

	log.Printf("successfully get user by email %s", email)
	return userRow.ToEntity(), true, nil
}

func (r *UsersRepository) UpdateUser(ctx context.Context, user entity.User) error {
	userRow := pgEntity.NewUserRow().FromEntity(user)
	if err := r.Update(ctx, userRow, userRow.ConditionUserUidEqual()); err != nil {
		log.Printf("failed to update user with uid %s. error: %v", user.Uid, err)
		return errors.Wrapf(err, "failed to update user with uid %s", user.Uid)
	}

	log.Printf("successfully update user with uid %s", user.Uid.String())
	return nil
}

func (r *UsersRepository) DeleteUser(ctx context.Context, uid uuid.UUID) error {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Uid: uid})

	if err := r.Delete(ctx, userRow, squirrel.And{userRow.ConditionUserUidEqual(), userRow.ConditionNotSuperAdmin()}); err != nil {
		log.Printf("failed to delete user with uid %s. error: %v", uid, err)
		return errors.WithStack(err)
	}

	log.Printf("successfully delete user with uid: %s", uid)
	return nil
}
