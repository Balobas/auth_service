package repositoryUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *UsersRepository) CreateUserWithPermissions(ctx context.Context, user entity.User) error {
	if err := r.WithTx(ctx, func(ctx context.Context) error {

		for _, row := range []repositoryPostgres.Row{
			pgEntity.NewUserRow().FromEntity(user),
			pgEntity.NewUserPermissionsRow().FromEntity(user),
		} {
			if err := r.Create(ctx, row); err != nil {
				return errors.Wrapf(err, "failed to create row in table %s", row.Table())
			}
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	log.Printf("successfuly create user uid: %s, name: %s, email: %s, permissions: %v\n", user.Uid, user.Name, user.Email, user.Permissions)
	return nil
}

func (r *UsersRepository) GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Uid: uid})

	if err := r.GetOne(ctx, userRow, userRow.ConditionUserUidEqual()); err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	permissionsRow := pgEntity.NewUserPermissionsRow()
	if err := r.GetOne(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	user := userRow.ToEntity()
	permissionsRow.ToEntity(&user)

	return user, nil
}

func (r *UsersRepository) UpdateUserWithPermissions(ctx context.Context, user entity.User) error {
	if err := r.WithTx(ctx, func(ctx context.Context) error {
		userRow := pgEntity.NewUserRow().FromEntity(user)
		permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(user)

		if err := r.Update(ctx, userRow, userRow.ConditionUserUidEqual()); err != nil {
			return errors.Wrapf(err, "failed to update user with uid %s", user.Uid)
		}

		if err := r.Update(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
			return errors.Wrapf(err, "failed to update user permissions (user uid: %s)", user.Uid)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	log.Printf("update user with uid: %d", user.Uid)
	return nil
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
