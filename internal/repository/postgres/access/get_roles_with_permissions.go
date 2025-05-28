package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
)

func (r *AccessRepository) GetRolesWithPermissions(ctx context.Context, limit int64, offset int64) ([]entity.RoleWithPermissions, error) {
	log.Printf("accessRepository.GetRolesWithPermissions: limit %d, offset %d", limit, offset)

	//TODO::
	return nil, nil
}
