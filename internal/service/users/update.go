package usersService

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
)

func (s *UsersService) Update(ctx context.Context, user entity.User) error {
	log.Printf("service.Update id: %d", user.Id)
	return s.repo.UpdateUser(ctx, user)
}