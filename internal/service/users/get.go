package usersService

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
)

func (s *UsersService) Get(ctx context.Context, id int64) (entity.User, error) {
	log.Printf("sevice.Get id: %d", id)
	return s.repo.GetUser(ctx, id)
}