package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (s *AuthServerGrpc) GetRolePermissions(ctx context.Context, req *auth_v1.GetRolePermissionsParams) (*auth_v1.Permissions, error) {
	log.Printf("authServerGrpc.GetRolePermissions: role %s", req.GetRole())

	userInfo := userInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetRolePermissions: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	permissions, err := s.ucAccess.GetRolePermissions(ctx, req.GetRole())
	if err != nil {
		log.Printf("authServerGrpc.GetRolePermissions: failed to get role permissions: %v", err)
		return nil, errors.WithStack(err)
	}

	permissionsPb := make([]*auth_v1.Permission, len(permissions))
	for i, permission := range permissions {
		permissionsPb[i] = &auth_v1.Permission{Key: permission.Key, Description: permission.Description}
	}

	return &auth_v1.Permissions{Permissions: permissionsPb}, nil
}
