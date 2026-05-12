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

func (s *AuthServerGrpc) DeleteRole(ctx context.Context, req *auth_v1.Role) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.DeleteRole: role %s", req.GetName())

	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.DeleteRole: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucAccess.DeleteRole(ctx, req.GetName()); err != nil {
		log.Printf("authServerGrpc.DeleteRole: failed to delete role %s: %v", req.GetName(), err)
		return nil, errors.WithStack(err)
	}

	return &emptypb.Empty{}, nil
}
