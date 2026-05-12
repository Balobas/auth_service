package deliveryGrpcAuth

import (
	"context"

	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthServerGrpc) GetUserAuthorizedDevices(ctx context.Context, req *auth_v1.UUID) (*auth_v1.AuthorizedDevices, error) {
	userUid, err := uuid.FromString(req.GetUid())
	if err != nil {
		return nil, errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}
	devices, err := s.ucDevices.GetUserAuthorizedDevices(ctx, userUid)
	if err != nil {
		return nil, err
	}

	res := make([]*auth_v1.AuthorizedDevice, len(devices))
	for i := 0; i < len(devices); i++ {
		res[i] = &auth_v1.AuthorizedDevice{
			Uid:          devices[i].Uid.String(),
			UserUid:      devices[i].MaintainerUid.String(),
			Name:         devices[i].Name,
			Agent:        devices[i].Agent,
			Language:     devices[i].Language,
			CreatedAt:    timestamppb.New(devices[i].CreatedAt),
			AuthorizedAt: timestamppb.New(devices[i].AuthorizedAt),
		}
	}
	return &auth_v1.AuthorizedDevices{
		Devices: res,
	}, nil
}
