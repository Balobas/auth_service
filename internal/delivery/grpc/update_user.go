package deliveryGrpc

import (
	"context"
	"errors"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
)

func (s *AuthServerGrpc) UpdateUser(ctx context.Context, req *auth_v1.UpdateUserRequest) (*auth_v1.JwtResponse, error) {
	userInfo := userInfoFromContext(ctx)

	userUid := uuid.FromStringOrNil(req.GetUid())

	if !uuid.Equal(userUid, userInfo.UserUid) {
		return nil, errors.New("permissions denied")
	}

	access, refresh, err := s.ucAuth.UpdateUserCreds(
		ctx,
		entity.User{
			Uid:   userUid,
			Email: req.GetEmail(),
			Role:  userInfo.Role,
		},
		req.GetPassword(),
	)
	if err != nil {
		return nil, err
	}

	return &auth_v1.JwtResponse{
		AccessJwt:  access,
		RefreshJwt: refresh,
	}, nil
}
