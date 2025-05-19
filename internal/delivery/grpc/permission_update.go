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

func (s *AuthServerGrpc) UpdatePermission(ctx context.Context, req *auth_v1.Permission) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.UpdatePermission: key %s description %s", req.GetKey(), req.GetDescription())

	userInfo := userInfoFromContext(ctx)

	if userInfo.Role != entity.UserRoleAdmin {
		log.Printf("authServerGrpc.UpdatePermission: user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucPermissions.UpdatePermission(ctx, entity.Permission{
		Key:         req.GetKey(),
		Description: req.GetDescription(),
	}); err != nil {
		log.Printf("authServerGrpc.UpdatePermission: failed to update permission %s: %v", req.GetKey(), err)
		return nil, errors.WithStack(err)
	}
	return nil, nil
}
