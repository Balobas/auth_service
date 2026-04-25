package deliveryGrpc

import (
	"context"

	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthServerGrpc) GetUsersAuthorizedDevices(ctx context.Context, req *auth_v1.UUIDs) (*auth_v1.AuthorizedDevices, error) {
	if req == nil || len(req.Uids) == 0 {
		return nil, errors.Wrap(serviceErrors.ErrBadRequest, "empty users uids")
	}

	usersUids, err := convertUidsFromStrings(req.Uids)
	if err != nil {
		return nil, errors.Wrap(serviceErrors.ErrBadRequest, "invalid users uids")
	}

	devices, err := s.ucDevices.GetUsersAuthorizedDevices(ctx, usersUids...)
	if err != nil {
		return nil, err
	}

	res := make([]*auth_v1.AuthorizedDevice, len(devices))
	for i := 0; i < len(devices); i++ {
		res[i] = &auth_v1.AuthorizedDevice{
			Uid:          devices[i].Uid.String(),
			UserUid:      devices[i].UserUid.String(),
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

func convertUidsFromStrings(uids []string) ([]uuid.UUID, error) {
	res := make([]uuid.UUID, len(uids))
	var err error

	for i := 0; i < len(uids); i++ {
		res[i], err = uuid.FromString(uids[i])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
