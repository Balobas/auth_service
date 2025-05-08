package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	"github.com/pkg/errors"
)

func (s *AuthServerGrpc) GetPermissions(ctx context.Context, req *auth_v1.GetPermissionsParams) (*auth_v1.Permissions, error) {
	log.Printf("authServerGrpc.GetPermissions: permissionPattern %s", req.GetPermissionPattern())

	userInfo := userInfoFromContext(ctx)

	if userInfo.Role != entity.UserRoleAdmin {
		log.Printf("authServerGrpc.GetPermissions: user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.New("permission denied")
	}

	perms, err := s.ucPermissions.GetPermissions(ctx, req.GetPermissionPattern())
	if err != nil {
		log.Printf("authServerGrpc.GetPermissions: failed to get permissions: %v", err)
		return nil, errors.WithStack(err)
	}

	if len(perms) == 0 {
		return &auth_v1.Permissions{Permissions: []*auth_v1.Permission{}}, nil
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
