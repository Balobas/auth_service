package usersService

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
)

func (s *UsersService) Create(ctx context.Context, user entity.User) (int64, error) {
	log.Printf("service.create user name: %v, email: %v", user.Name, user.Email)

	return s.repo.CreateUser(ctx, user)
}
