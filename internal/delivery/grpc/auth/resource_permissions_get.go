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

func (s *AuthServerGrpc) GetResourcePermission(ctx context.Context, req *auth_v1.GetResourcePermissionParams) (*auth_v1.Permission, error) {
	log.Printf("authServerGrpc.GetResourcePermission: resource %s method %s", req.GetUri(), req.GetMethod())

	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetResourcePermission: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	permission, err := s.ucAccess.GetResourcePermission(ctx, req.GetUri(), req.GetMethod())
	if err != nil {
		log.Printf("authServerGrpc.GetResourcePermission: failed to get resource %s method %s permission: %v", req.GetUri(), req.GetMethod(), err)
		return nil, errors.WithStack(err)
	}

	return &auth_v1.Permission{Key: permission}, nil
}
