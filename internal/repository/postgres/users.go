package repositoryPostgres

import (
	"context"
	"log"

	"github.com/balobas/auth_service_bln/internal/entity"
	pgEntity "github.com/balobas/auth_service_bln/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *Repository) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	userRow := pgEntity.NewUserRow().FromEntity(user)

	if err := r.create(ctx, userRow); err != nil {
		return 0, errors.WithStack(err)
	}

	log.Printf("successfuly create user id: %d, name: %s, email: %s \n", userRow.Id, user.Name, user.Email)
	return userRow.ToEntity().Id, nil
}

func (r *Repository) GetUser(ctx context.Context, id int64) (entity.User, error) {
	userRow := pgEntity.NewUserRow().FromEntity(entity.User{Id: id})

	if err := r.get(ctx, userRow); err != nil {
		return entity.User{}, errors.WithStack(err)
	}

	return userRow.ToEntity(), nil
}

func (r *Repository) UpdateUser(ctx context.Context, userParams entity.User) error {
	userRow := pgEntity.NewUserRow().FromEntity(userParams)

	if err := r.update(ctx, userRow); err != nil {
		return errors.WithStack(err)
	}

	log.Printf("update user with uid: %d", userParams.Id)
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {

	if err := r.delete(ctx, pgEntity.NewUserRow().FromEntity(entity.User{Id: id})); err != nil {
		return errors.WithStack(err)
	}

	log.Printf("delete user with uid: %d", id)
	return nil
}
