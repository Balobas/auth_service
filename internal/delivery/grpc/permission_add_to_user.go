package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) AddPermissionToUser(ctx context.Context, req *auth_v1.UserPermissionParams) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.AddPermissionToUser: user %s permission %s", req.GetUserUid(), req.GetPermission())

	userInfo := userInfoFromContext(ctx)

	if userInfo.Role != entity.UserRoleAdmin {
		log.Printf("authServerGrpc.AddPermissionToUser: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucPermissions.AddPermissionToUser(ctx, uuid.FromStringOrNil(req.GetUserUid()), req.GetPermission()); err != nil {
		log.Printf("authServerGrpc.AddPermissionToUser: failed to add permission %s to user %s: %v", req.GetPermission(), req.GetUserUid(), err)
		return nil, errors.WithStack(err)
	}
	return nil, nil
}
