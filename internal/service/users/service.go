package usersService

import (
	"context"

	"github.com/balobas/auth_service_bln/internal/entity"
)

type UsersService struct {
	repo UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	GetUser(ctx context.Context, id int64) (entity.User, error)
	UpdateUser(ctx context.Context, userParams entity.User) error
	DeleteUser(ctx context.Context, id int64) error
}

func New(repo UsersRepository) *UsersService {
	return &UsersService{
		repo: repo,
	}
}
