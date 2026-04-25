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

func (r *Repository) CreateUserAuthorizedDevice(ctx context.Context, device entity.UserAuthorizedDevice) error {
	return r.Create(ctx, pgEntity.NewDeviceRow().FromEntity(device))
}

func (r *Repository) UpdateUserAuthorizedDevice(ctx context.Context, device entity.UserAuthorizedDevice) error {
	row := pgEntity.NewDeviceRow().FromEntity(device)
	return r.Update(ctx, row, row.ConditionUserUidAndDeviceUidEqual())
}

func (r *Repository) GetUserAuthorizedDevice(ctx context.Context, userUid uuid.UUID, deviceUid uuid.UUID) (entity.UserAuthorizedDevice, bool, error) {
	row := pgEntity.NewDeviceRow().FromEntity(entity.UserAuthorizedDevice{UserDevice: entity.UserDevice{Uid: deviceUid, UserUid: userUid}})
	if err := r.GetOne(ctx, row, row.ConditionUserUidAndDeviceUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserAuthorizedDevice{}, false, nil
		}
		return entity.UserAuthorizedDevice{}, false, err
	}
	return row.ToEntity(), true, nil
}

func (r *Repository) GetUserAuthorizedDevices(ctx context.Context, userUid uuid.UUID) ([]entity.UserAuthorizedDevice, error) {
	row := pgEntity.NewDeviceRow().FromEntity(entity.UserAuthorizedDevice{UserDevice: entity.UserDevice{UserUid: userUid}})

	rows := pgEntity.NewDevicesRows()
	if err := r.GetSome(ctx, row, rows,
		sq.And{
			row.ConditionUserUidEqual(),
			row.ConditionUnauthorizedIsNull(),
		},
	); err != nil {
		return nil, err
	}

	return rows.ToEntities(), nil
}

func (r *Repository) GetUsersAuthorizedDevices(ctx context.Context, usersUids ...uuid.UUID) ([]entity.UserAuthorizedDevice, error) {
	row := pgEntity.NewDeviceRow()

	rows := pgEntity.NewDevicesRows()
	if err := r.GetSome(ctx, row, rows,
		sq.And{
			row.ConditionUserUidIn(usersUids),
			row.ConditionUnauthorizedIsNull(),
		},
	); err != nil {
		return nil, err
	}

	return rows.ToEntities(), nil
}
