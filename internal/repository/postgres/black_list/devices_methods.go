package blacklistRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *BlacklistRepository) CreateBlackListDevice(ctx context.Context, blackListDevice entity.BlackListDevice) error {
	row := pgEntity.NewDevicesBlackListRow().FromBlackListDeviceEntity(blackListDevice)

	if err := r.Create(ctx, row); err != nil {
		return errors.Wrapf(err, "failed to insert device with uid %s to black list", blackListDevice.Device.Uid)
	}

	return nil
}

func (r *BlacklistRepository) GetBlackListDevice(ctx context.Context, deviceUid uuid.UUID) (entity.BlackListElement, error) {
	row := pgEntity.NewDevicesBlackListRow().FromBlackListEntity(entity.BlackListElement{Uid: deviceUid})

	if err := r.GetOne(ctx, row, row.ConditionDeviceUidEqual()); err != nil {
		return entity.BlackListElement{}, errors.Wrapf(err, "failed to get device with uid %s", deviceUid)
	}

	return row.ToBlackListEntity(), nil
}

func (r *BlacklistRepository) UpdateBlackListDevice(ctx context.Context, elem entity.BlackListElement) error {
	row := pgEntity.NewDevicesBlackListRow().FromBlackListEntity(elem)

	if err := r.Update(ctx, row, row.ConditionDeviceUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to update black list device with uid %s", elem.Uid)
	}

	return nil
}

func (r *BlacklistRepository) DeleteDeviceFromBlackList(ctx context.Context, deviceUid uuid.UUID) error {
	row := pgEntity.NewDevicesBlackListRow().FromBlackListEntity(entity.BlackListElement{Uid: deviceUid})

	if err := r.Delete(ctx, row, row.ConditionDeviceUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete device with uid %s from black list", deviceUid)
	}

	return nil
}

func (r *BlacklistRepository) GetBlackListDevices(ctx context.Context, limit, offset int64) ([]entity.BlackListDevice, error) {
	deviceRow := &pgEntity.DeviceRow{}
	blackListDeviceRow := &pgEntity.DevicesBlackListRow{}

	devicePreffix := "d"
	blackListDevicePreffix := "b"

	stmt := fmt.Sprintf(
		"SELECT %s, %s FROM %s AS %s JOIN %s AS %s ON %s=%s LIMIT %d OFFSET %d",
		strings.Join(makePreffix(devicePreffix, deviceRow.Columns()), ","),
		strings.Join(makePreffix(blackListDevicePreffix, blackListDeviceRow.Columns()), ","),
		deviceRow.Table(), devicePreffix,
		blackListDeviceRow.Table(), blackListDevicePreffix,
		devicePreffix+"."+deviceRow.IdColumnName(), blackListDevicePreffix+"."+blackListDeviceRow.IdColumnName(),
		limit, offset,
	)

	rows, err := r.DB().Query(ctx, stmt, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]entity.BlackListDevice, 0, limit)

	for rows.Next() {
		if err := rows.Scan(append(deviceRow.ValuesForScan(), blackListDeviceRow.ValuesForScan()...)...); err != nil {
			return nil, errors.Wrap(err, "failed to scan")
		}

		res = append(res, entity.BlackListDevice{
			Device: deviceRow.ToEntity(),
			Reason: blackListDeviceRow.ToBlackListEntity().Reason,
		})
	}

	return res, nil
}
