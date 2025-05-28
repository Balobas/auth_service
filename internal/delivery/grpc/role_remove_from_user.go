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

func (s *AuthServerGrpc) RemoveRoleFromUser(ctx context.Context, req *auth_v1.UserRoleParams) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.RemoveRoleFromUser: user %s role %s", req.GetUserUid(), req.GetRole())

	userInfo := userInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.RemoveRoleFromUser: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucAccess.RemoveRoleFromUser(ctx, uuid.FromStringOrNil(req.GetUserUid()), req.GetRole()); err != nil {
		log.Printf("authServerGrpc.RemoveRoleFromUser: failed to remove role %s from user %s: %v", req.GetRole(), req.GetUserUid(), err)
		return nil, errors.WithStack(err)
	}
	return nil, nil
}
