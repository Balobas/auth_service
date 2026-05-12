package devicesRepository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

func (r *Repository) CreateAuthorizedDevice(ctx context.Context, device entity.AuthorizedDevice) error {
	return r.Create(ctx, pgEntity.NewDeviceRow().FromEntity(device))
}

func (r *Repository) UpdateAuthorizedDevice(ctx context.Context, device entity.AuthorizedDevice) error {
	row := pgEntity.NewDeviceRow().FromEntity(device)
	return r.Update(ctx, row, row.ConditionMaintainerUidAndDeviceUidEqual())
}

func (r *Repository) GetAuthorizedDevice(ctx context.Context, maintainerUid uuid.UUID, deviceUid uuid.UUID) (entity.AuthorizedDevice, bool, error) {
	row := pgEntity.NewDeviceRow().FromEntity(entity.AuthorizedDevice{Device: entity.Device{Uid: deviceUid, MaintainerUid: maintainerUid}})
	if err := r.GetOne(ctx, row, row.ConditionMaintainerUidAndDeviceUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.AuthorizedDevice{}, false, nil
		}
		return entity.AuthorizedDevice{}, false, err
	}
	return row.ToEntity(), true, nil
}

func (r *Repository) GetUserAuthorizedDevices(ctx context.Context, userUid uuid.UUID) ([]entity.AuthorizedDevice, error) {
	row := pgEntity.NewDeviceRow().FromEntity(entity.AuthorizedDevice{Device: entity.Device{MaintainerUid: userUid}})

	rows := pgEntity.NewDevicesRows()
	if err := r.GetSome(ctx, row, rows,
		sq.And{
			row.ConditionMaintainerUidEqual(),
			row.ConditionUnauthorizedIsNull(),
			row.ConditionTypeEqual(entity.DeviceTypeUserDevice),
		},
	); err != nil {
		return nil, err
	}

	return rows.ToEntities(), nil
}

func (r *Repository) GetUsersAuthorizedDevices(ctx context.Context, usersUids ...uuid.UUID) ([]entity.AuthorizedDevice, error) {
	row := pgEntity.NewDeviceRow()

	rows := pgEntity.NewDevicesRows()
	if err := r.GetSome(ctx, row, rows,
		sq.And{
			row.ConditionMaintainerUidIn(usersUids),
			row.ConditionUnauthorizedIsNull(),
			row.ConditionTypeEqual(entity.DeviceTypeUserDevice),
		},
	); err != nil {
		return nil, err
	}

	return rows.ToEntities(), nil
}
