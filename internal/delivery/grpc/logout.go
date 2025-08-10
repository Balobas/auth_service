gpackage deliveryGrpc

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

func (s *AuthServerGrpc) Logout(ctx context.Context, req *auth_v1.LogoutRequest) (*emptypb.Empty, error) {
	userInfo := userInfoFromContext(ctx)
	log.Printf("authServerGrpc.Logout: user uid %s, caller %s", req.GetUid(), userInfo.UserUid)

	userUid := uuid.FromStringOrNil(req.GetUid())
	if uuid.Equal(userUid, uuid.UUID{}) {
		userUid = userInfo.UserUid
	}

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) && !uuid.Equal(userInfo.UserUid, userUid) {
		log.Printf("authServerGrpc.Logout: user %s cant logout user %s", userInfo.UserUid, userUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "permissions denied")
	}

	if err := s.ucAuth.Logout(ctx, userUid); err != nil {
		log.Printf("authServerGrpc.Logout: failed to logout user %s: %v", userUid, err)
		return nil, err
	}
	return nil, nil
}
