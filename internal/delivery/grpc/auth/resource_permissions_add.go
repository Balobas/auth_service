package deliveryGrpcAuth

import (
	"context"
	"log"

	deliveryGrpcInterceptors "github.com/balobas/auth_service/internal/delivery/grpc/interceptors"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) AddResourcePermission(ctx context.Context, req *auth_v1.ResourcePermissionParams) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.AddResourcePermission: resource %s method %s permission %s", req.GetUri(), req.GetMethod(), req.GetPermission())

	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.AddResourcePermission: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucAccess.AddResourcePermission(ctx, req.GetUri(), req.GetMethod(), req.GetPermission()); err != nil {
		log.Printf("authServerGrpc.AddResourcePermission: failed to add resource %s method %s permission %s: %v", req.GetUri(), req.GetMethod(), req.GetPermission(), err)
		return nil, errors.WithStack(err)
	}

	return &emptypb.Empty{}, nil
}
