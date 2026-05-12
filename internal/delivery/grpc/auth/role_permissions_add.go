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

func (s *AuthServerGrpc) AddPermissionsToRole(ctx context.Context, req *auth_v1.RolePermissionsParams) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.AddPermissionsToRole: role %s permissions %v", req.GetRole(), req.GetPermissions())

	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.AddPermissionsToRole: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucAccess.AddPermissionsToRole(ctx, req.GetRole(), req.GetPermissions()); err != nil {
		log.Printf("authServerGrpc.AddPermissionsToRole: failed to add permissions %v to role %s: %v", req.GetPermissions(), req.GetRole(), err)
		return nil, errors.WithStack(err)
	}

	return &emptypb.Empty{}, nil
}
