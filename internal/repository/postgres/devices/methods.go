package devicesRepository

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *DevicesRepository) CreateDevice(ctx context.Context, device entity.UserDevice) error {
	deviceRow := pgEntity.NewDeviceRow().FromEntity(device)
	if err := r.Create(ctx, deviceRow); err != nil {
		return errors.Wrapf(err, "failed to create device with uid %s and user uid %s", device.Uid, device.UserUid)
	}
	return nil
}

func (r *DevicesRepository) GetDeviceByUid(ctx context.Context, uid uuid.UUID) (entity.UserDevice, error) {
	deviceRow := pgEntity.NewDeviceRow().FromEntity(entity.UserDevice{Uid: uid})

	if err := r.Get(ctx, deviceRow); err != nil {
		return entity.UserDevice{}, errors.Wrapf(err, "failed to get device by uid %s", uid)
	}

	return deviceRow.ToEntity(), nil
}

func (r *DevicesRepository) GetUserDevices(ctx context.Context, userUid uuid.UUID) ([]entity.UserDevice, error) {
	deviceRow := pgEntity.NewDeviceRow().FromEntity(entity.UserDevice{UserUid: userUid})
	resultRows := pgEntity.NewDeviceRows()

	if err := r.GetByCondition(ctx, deviceRow, resultRows, deviceRow.GetByUserUidCondition()); err != nil {
		return nil, errors.WithStack(err)
	}

	return resultRows.ToEntities(), nil
}

func (r *DevicesRepository) DeleteDevice(ctx context.Context, uid uuid.UUID) error {
	deviceRow := pgEntity.NewDeviceRow().FromEntity(entity.UserDevice{Uid: uid})
	if err := r.Delete(ctx, deviceRow); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
