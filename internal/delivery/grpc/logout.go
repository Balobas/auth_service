package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) Logout(ctx context.Context, req *auth_v1.LogoutRequest) (*emptypb.Empty, error) {
	userInfo := userInfoFromContext(ctx)
	log.Printf("authServerGrpc.Logout: user uid %s device uid %s", userInfo.UserUid, req.GetDeviceUid())

	if err := s.ucAuth.Logout(ctx, userInfo.UserUid, uuid.FromStringOrNil(req.GetDeviceUid())); err != nil {
		log.Printf("authServerGrpc.Logout: failed to logout user %s from device %s: %v", userInfo.UserUid, req.GetDeviceUid(), err)
		return nil, err
	}
	return nil, nil
}
