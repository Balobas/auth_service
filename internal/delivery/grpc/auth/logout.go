package deliveryGrpcAuth

import (
	"context"
	"log"

	deliveryGrpcInterceptors "github.com/balobas/auth_service/internal/delivery/grpc/interceptors"
	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) Logout(ctx context.Context, req *auth_v1.LogoutRequest) (*emptypb.Empty, error) {
	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)
	log.Printf("authServerGrpc.Logout: user uid %s device uid %s", userInfo.UserUid, req.GetDeviceUid())

	var logoutDeviceUid uuid.UUID
	deviceUidFromRequest := uuid.FromStringOrNil(req.GetDeviceUid())
	if !uuid.Equal(deviceUidFromRequest, uuid.UUID{}) {
		logoutDeviceUid = deviceUidFromRequest
	} else {
		// Если в запросе не указан явно uid девайса -> берем из токена
		logoutDeviceUid = userInfo.DeviceUid
	}

	if err := s.ucAuth.Logout(ctx, userInfo.UserUid, logoutDeviceUid); err != nil {
		log.Printf("authServerGrpc.Logout: failed to logout user %s from device %s: %v", userInfo.UserUid, req.GetDeviceUid(), err)
		return nil, err
	}
	return nil, nil
}
