package deliveryGrpc

import (
	"context"
	"errors"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) UpdateUser(ctx context.Context, req *auth_v1.UpdateUserRequest) (*emptypb.Empty, error) {
	userInfo := userInfoFromContext(ctx)

	userUid := uuid.FromStringOrNil(req.GetUid())

	if !uuid.Equal(userUid, userInfo.UserUid) {
		return nil, errors.New("permissions denied")
	}

	if err := s.ucUsers.UpdateUser(
		ctx,
		entity.User{
			Uid:   userUid,
			Email: req.GetEmail(),
		},
		req.GetPassword(),
	); err != nil {
		return nil, err
	}
	return nil, nil
}
