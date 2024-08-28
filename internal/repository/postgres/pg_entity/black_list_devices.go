package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const devicesBlackListTableName = "black_list_devices"

type DevicesBlackListRow struct {
	DeviceUid pgtype.UUID
	Reason    string
}

func NewDevicesBlackListRow() *DevicesBlackListRow {
	return &DevicesBlackListRow{}
}

var (
	devicesBlackListColumns = []string{
		"device_uid",
		"reason",
	}
)

func (u *DevicesBlackListRow) FromBlackListEntity(elem entity.BlackListElement) *DevicesBlackListRow {
	u.DeviceUid = pgtype.UUID{
		Bytes:  elem.Uid,
		Status: pgtype.Present,
	}
	u.Reason = elem.Reason
	return u
}

func (u *DevicesBlackListRow) ToBlackListEntity() entity.BlackListElement {
	return entity.BlackListElement{
		Uid:    u.DeviceUid.Bytes,
		Reason: u.Reason,
	}
}

func (u *DevicesBlackListRow) FromBlackListDeviceEntity(device entity.BlackListDevice) *DevicesBlackListRow {
	u.DeviceUid = pgtype.UUID{
		Bytes:  device.Device.Uid,
		Status: pgtype.Present,
	}
	u.Reason = device.Reason
	return u
}

func (u *DevicesBlackListRow) IdColumnName() string {
	return "device_uid"
}

func (u *DevicesBlackListRow) Values() []interface{} {
	return []interface{}{
		u.DeviceUid,
		u.Reason,
	}
}

func (u *DevicesBlackListRow) Columns() []string {
	return devicesBlackListColumns
}

func (u *DevicesBlackListRow) Table() string {
	return devicesBlackListTableName
}

func (u *DevicesBlackListRow) Scan(row pgx.Row) error {
	return row.Scan(&u.DeviceUid, &u.Reason)
}

func (u *DevicesBlackListRow) ColumnsForUpdate() []string {
	return []string{"reason"}
}

func (u *DevicesBlackListRow) ValuesForUpdate() []interface{} {
	return []interface{}{u.Reason}
}

func (u *DevicesBlackListRow) ValuesForScan() []interface{} {
	return []interface{}{
		&u.DeviceUid,
		&u.Reason,
	}
}

func (u *DevicesBlackListRow) ConditionDeviceUidEqual() sq.Eq {
	return sq.Eq{
		"device_uid": u.DeviceUid,
	}
}
