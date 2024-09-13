package deliveryGrpc

import (
	"context"
	"errors"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) Delete(ctx context.Context, req *auth_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	userInfo := userInfoFromContext(ctx)

	userUid := uuid.FromStringOrNil(req.GetUid())

	if userInfo.Role != entity.UserRoleAdmin && !uuid.Equal(userUid, userInfo.UserUid) {
		return nil, errors.New("access denied")
	}

	if err := s.ucUsers.DeleteUser(ctx, userUid); err != nil {
		return nil, err
	}
	return nil, nil
}
