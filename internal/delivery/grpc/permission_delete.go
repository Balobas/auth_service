package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) DeletePermission(ctx context.Context, req *auth_v1.Permission) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.DeletePermission: key %s description %s", req.GetKey(), req.GetDescription())

	userInfo := userInfoFromContext(ctx)

	if userInfo.Role != entity.UserRoleAdmin {
		log.Printf("authServerGrpc.DeletePermission: user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucPermissions.DeletePermission(ctx, req.GetKey()); err != nil {
		log.Printf("authServerGrpc.DeletePermission: failed to delete permission %s: %v", req.GetKey(), err)
		return nil, errors.WithStack(err)
	}
	return nil, nil
}
