package deliveryGrpcAuth

import (
	"context"
	"log"

	deliveryGrpcInterceptors "github.com/balobas/auth_service/internal/delivery/grpc/interceptors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) LogoutAllOtherDevices(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)
	log.Printf("authServerGrpc.LogoutAllOtherDevices: user uid %s device uid %s", userInfo.UserUid, userInfo.DeviceUid)

	if err := s.ucAuth.LogoutAllOtherDevices(ctx, userInfo.UserUid, userInfo.DeviceUid); err != nil {
		log.Printf("authServerGrpc.LogoutAllOtherDevices: failed to logout user %s from all devices exclude current device %s: %v", userInfo.UserUid, userInfo.DeviceUid, err)
		return nil, err
	}
	return nil, nil
}
