package deliveryGrpcAuth

import (
	"context"
	"log"

	deliveryGrpcInterceptors "github.com/balobas/auth_service/internal/delivery/grpc/interceptors"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (s *AuthServerGrpc) GetResourcesPermissions(ctx context.Context, req *auth_v1.GetResourcesPermissionsParams) (*auth_v1.ResourcesPermissions, error) {
	log.Printf("authServerGrpc.GetResourcesPermissions: limit %d offset %d", req.GetLimit(), req.GetOffset())

	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetResourcesPermissions: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	permissions, err := s.ucAccess.GetResourcesPermissions(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		log.Printf("authServerGrpc.GetResourcesPermissions: failed to get resources permissions: %v", err)
		return nil, errors.WithStack(err)
	}

	permissionsPb := make([]*auth_v1.ResourcePermissionParams, len(permissions))
	for i, permission := range permissions {
		permissionsPb[i] = &auth_v1.ResourcePermissionParams{Uri: permission.URI, Method: permission.Method, Permission: permission.Permission}
	}

	return &auth_v1.ResourcesPermissions{ResourcesPermissions: permissionsPb}, nil
}
