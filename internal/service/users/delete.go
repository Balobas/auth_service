package usersService

import (
	"context"
	"log"
)

func (s *UsersService) Delete(ctx context.Context, id int64) error {
	log.Printf("service.Delete id: %d", id)
	return s.repo.DeleteUser(ctx, id)
}
