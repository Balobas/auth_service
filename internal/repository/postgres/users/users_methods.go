package repositoryUsers

import (
	"context"
	"log"

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

	log.Printf("successfuly create user uid: %s, email: %s, permissions: %v\n", user.Uid, user.Email, user.Permissions)
	return nil
}

func (r *UsersRepository) GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, bool, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Uid: uid})

	if err := r.GetOne(ctx, userRow, userRow.ConditionUserUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, errors.WithStack(err)
	}

	return userRow.ToEntity(), true, nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (entity.User, bool, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Email: email})

	if err := r.GetOne(ctx, userRow, userRow.ConditionEmailEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, errors.WithStack(err)
	}

	return userRow.ToEntity(), true, nil
}

func (r *UsersRepository) UpdateUser(ctx context.Context, user entity.User) error {
	userRow := pgEntity.NewUserRow().FromEntity(user)
	if err := r.Update(ctx, userRow, userRow.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to update user with uid %s", user.Uid)
	}
	return nil
}

func (r *UsersRepository) DeleteUser(ctx context.Context, uid uuid.UUID) error {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Uid: uid})
	if err := r.Delete(ctx, userRow, userRow.ConditionUserUidEqual()); err != nil {
		return errors.WithStack(err)
	}

	log.Printf("delete user with uid: %s", uid)
	return nil
}
