package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (s *AuthServerGrpc) GetPermissions(ctx context.Context, req *auth_v1.GetPermissionsParams) (*auth_v1.Permissions, error) {
	log.Printf("authServerGrpc.GetPermissions: permissionPattern %s", req.GetPermissionPattern())

	userInfo := userInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetPermissions: user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	perms, err := s.ucAccess.GetPermissions(ctx, req.GetPermissionPattern())
	if err != nil {
		log.Printf("authServerGrpc.GetPermissions: failed to get permissions: %v", err)
		return nil, errors.WithStack(err)
	}

	res := make([]*auth_v1.Permission, len(perms))
	for i := 0; i < len(perms); i++ {
		res[i] = &auth_v1.Permission{
			Key:         perms[i].Key,
			Description: perms[i].Description,
		}
	}

	return &auth_v1.Permissions{
		Permissions: res,
	}, nil
}
