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

func (s *AuthServerGrpc) CreatePermission(ctx context.Context, req *auth_v1.Permission) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.CreatePermission: key %s description %s", req.GetKey(), req.GetDescription())

	userInfo := userInfoFromContext(ctx)

	if userInfo.Role != entity.UserRoleAdmin {
		log.Printf("authServerGrpc.CreatePermission: user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucPermissions.CreatePermission(ctx, entity.Permission{
		Key:         req.GetKey(),
		Description: req.GetDescription(),
	}); err != nil {
		log.Printf("authServerGrpc.CreatePermission: failed to create permission %s: %v", req.GetKey(), err)
		return nil, errors.WithStack(err)
	}
	return nil, nil
}
