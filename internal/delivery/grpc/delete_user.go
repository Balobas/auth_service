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

func (s *AuthServerGrpc) DeleteUser(ctx context.Context, req *auth_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	userInfo := userInfoFromContext(ctx)
	log.Printf("authServerGrpc.DeleteUser: user uid %s, caller %s", req.GetUid(), userInfo.UserUid)

	userUid := uuid.FromStringOrNil(req.GetUid())

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) && !uuid.Equal(userUid, userInfo.UserUid) {
		log.Printf("authServerGrpc.DeleteUser: user %s hasnt permissions to delete user %s", userInfo.UserUid, userUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "permissions denied")
	}

	if err := s.ucUsers.DeleteUser(ctx, userUid); err != nil {
		log.Printf("authServerGrpc.DeleteUser: failed to delete user %s: %v", userUid, err)
		return nil, err
	}
	return nil, nil
}
