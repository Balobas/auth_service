package repositoryUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/internal/entity/contract"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user entity.User) (err error) {
	if !r.DB().HasTxInCtx(ctx) {
		var (
			tx         contract.Transaction
			beginTxErr error
		)
		ctx, tx, beginTxErr = r.DB().BeginTxWithContext(ctx)
		if beginTxErr != nil {
			return errors.Wrap(beginTxErr, "failed to start transaction")
		}

		defer func() {
			err = HandleTxEnd(ctx, tx, err)
		}()
	}

	for _, row := range []repositoryPostgres.Row{
		pgEntity.NewUserRow().FromEntity(user),
		pgEntity.NewUserPermissionsRow().FromEntity(user),
	} {
		if err := r.Create(ctx, row); err != nil {
			return errors.Wrapf(err, "failed to create row in table %s", row.Table())
		}
	}

	log.Printf("successfuly create user uid: %s, name: %s, email: %s, permissions: %v\n", user.Uid, user.Name, user.Email, user.Permissions)
	return nil
}

func (r *UsersRepository) GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Uid: uid})

	if err := r.Get(ctx, userRow); err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	permissionsRow := pgEntity.NewUserPermissionsRow()
	if err := r.Get(ctx, permissionsRow); err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	user := userRow.ToEntity()
	permissionsRow.ToEntity(&user)

	return user, nil
}

func (r *UsersRepository) UpdateUser(ctx context.Context, user entity.User) (err error) {
	if !r.DB().HasTxInCtx(ctx) {
		var (
			tx         contract.Transaction
			beginTxErr error
		)
		ctx, tx, beginTxErr = r.DB().BeginTxWithContext(ctx)
		if beginTxErr != nil {
			return errors.Wrap(beginTxErr, "failed to start transaction")
		}

		defer func() {
			err = HandleTxEnd(ctx, tx, err)
		}()
	}

	for _, row := range []repositoryPostgres.Row{
		pgEntity.NewUserRow().FromEntity(user),
		pgEntity.NewUserPermissionsRow().FromEntity(user),
	} {
		if err := r.Update(ctx, row); err != nil {
			return errors.Wrapf(err, "failed to update table %s", row.Table())
		}
	}

	log.Printf("update user with uid: %d", user.Uid)
	return nil
}

func (r *UsersRepository) DeleteUser(ctx context.Context, uid uuid.UUID) error {

	if err := r.Delete(ctx, pgEntity.NewUserRow().FromEntity(entity.User{Uid: uid})); err != nil {
		return errors.WithStack(err)
	}

	log.Printf("delete user with uid: %s", uid)
	return nil
}
