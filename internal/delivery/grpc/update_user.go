package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (s *AuthServerGrpc) UpdateUser(ctx context.Context, req *auth_v1.UpdateUserRequest) (*auth_v1.JwtResponse, error) {
	userInfo := userInfoFromContext(ctx)
	log.Printf("authServerGrpc.UpdateUser: uid %s email %s caller %s", req.GetUid(), req.GetEmail(), userInfo.UserUid)

	userUid := uuid.FromStringOrNil(req.GetUid())

	if !uuid.Equal(userUid, uuid.UUID{}) && !uuid.Equal(userUid, userInfo.UserUid) {
		log.Printf("authServerGrpc.UpdateUser: failed to update another user: caller %s, user %s", userInfo.UserUid, userUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "cant update another user")
	}

	access, refresh, err := s.ucAuth.UpdateUserCreds(
		ctx,
		entity.User{
			Uid:   userUid,
			Email: req.GetEmail(),
		},
		req.GetPassword(),
		userInfo.DeviceUid,
	)
	if err != nil {
		log.Printf("authServerGrpc.UpdateUser: failed to update user %s: %v", userUid, err)
		return nil, err
	}

	return &auth_v1.JwtResponse{
		AccessJwt:  access,
		RefreshJwt: refresh,
	}, nil
}
