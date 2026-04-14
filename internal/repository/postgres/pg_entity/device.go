package pgEntity

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const devicesTableName = "users_devices"

var devicesTableColumns = []string{
	"uid",
	"user_uid",
	"name",
	"agent",
	"language",
	"created_at",
	"authorized_at",
	"unauthorized_at",
}

type DeviceRow struct {
	Uid            pgtype.UUID
	UserUid        pgtype.UUID
	Name           string
	Agent          string
	Language       string
	CreatedAt      pgtype.Timestamp
	AuthorizedAt   pgtype.Timestamp
	UnauthorizedAt pgtype.Timestamp
}

func NewDeviceRow() *DeviceRow {
	return &DeviceRow{}
}

func (s *DeviceRow) FromEntity(device entity.UserAuthorizedDevice) *DeviceRow {
	s.Uid = pgtype.UUID{
		Bytes: device.Uid,
		Valid: true,
	}
	s.UserUid = pgtype.UUID{
		Bytes: device.UserUid,
		Valid: true,
	}
	s.Name = device.Name
	s.Agent = device.Agent
	s.Language = device.Language

	if device.CreatedAt.Unix() == 0 || device.CreatedAt.IsZero() {
		s.CreatedAt = pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		}
	} else {
		s.CreatedAt = pgtype.Timestamp{
			Time:  device.CreatedAt.UTC(),
			Valid: true,
		}
	}

	if device.AuthorizedAt.Unix() == 0 || device.AuthorizedAt.IsZero() {
		s.AuthorizedAt = pgtype.Timestamp{}
	} else {
		s.AuthorizedAt = pgtype.Timestamp{
			Time:  device.AuthorizedAt.UTC(),
			Valid: true,
		}
	}

	if device.UnauthorizedAt == nil {
		s.UnauthorizedAt = pgtype.Timestamp{}
	} else {
		s.UnauthorizedAt = pgtype.Timestamp{
			Time:  device.AuthorizedAt.UTC(),
			Valid: true,
		}
	}
	return s
}

func (s *DeviceRow) ToEntity() entity.UserAuthorizedDevice {
	d := entity.UserAuthorizedDevice{
		UserDevice: entity.UserDevice{
			Uid:       s.Uid.Bytes,
			UserUid:   s.UserUid.Bytes,
			Name:      s.Name,
			Agent:     s.Agent,
			Language:  s.Language,
			CreatedAt: s.CreatedAt.Time,
		},
		AuthorizedAt: s.AuthorizedAt.Time,
	}
	if s.UnauthorizedAt.Valid {
		d.UnauthorizedAt = &s.UnauthorizedAt.Time
	}
	return d
}

func (s *DeviceRow) IdColumnName() string {
	return "uid"
}

func (s *DeviceRow) Values() []interface{} {
	return []interface{}{
		s.Uid,
		s.UserUid,
		s.Name,
		s.Agent,
		s.Language,
		s.CreatedAt,
		s.AuthorizedAt,
		s.UnauthorizedAt,
	}
}

func (s *DeviceRow) Columns() []string {
	return devicesTableColumns
}

func (s *DeviceRow) Table() string {
	return devicesTableName
}

func (s *DeviceRow) Scan(row pgx.Row) error {
	return row.Scan(
		&s.Uid,
		&s.UserUid,
		&s.Name,
		&s.Agent,
		&s.Language,
		&s.CreatedAt,
		&s.AuthorizedAt,
		&s.UnauthorizedAt,
	)
}

func (s *DeviceRow) ColumnsForUpdate() []string {
	return []string{
		"name",
		"agent",
		"language",
		"authorized_at",
		"unauthorized_at",
	}
}

func (s *DeviceRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		s.Name,
		s.Agent,
		s.Language,
		s.AuthorizedAt,
		s.UnauthorizedAt,
	}
}

func (s *DeviceRow) ConditionUidEqual() sq.Eq {
	return sq.Eq{
		"uid": s.Uid,
	}
}

func (s *DeviceRow) ConditionUserUidEqual() sq.Eq {
	return sq.Eq{
		"user_uid": s.UserUid,
	}
}

func (s *DeviceRow) ConditionUnauthorizedIsNull() sq.Eq {
	return sq.Eq{
		"unauthorized": pgtype.Timestamp{},
	}
}

func (s *DeviceRow) ConditionUserUidAndDeviceUidEqual() sq.Eq {
	return sq.Eq{
		"uid":      s.Uid,
		"user_uid": s.UserUid,
	}
}

type DevicesRows struct {
	devices []*DeviceRow
}

func NewDevicesRows() *DevicesRows {
	return &DevicesRows{}
}

func (s *DevicesRows) ScanAll(rows pgx.Rows) error {
	for rows.Next() {
		newRow := &DeviceRow{}

		if err := newRow.Scan(rows); err != nil {
			return err
		}
		s.devices = append(s.devices, newRow)
	}

	return nil
}

func (s *DevicesRows) ToEntities() []entity.UserAuthorizedDevice {
	if len(s.devices) == 0 {
		return nil
	}

	res := make([]entity.UserAuthorizedDevice, len(s.devices))

	for i := 0; i < len(s.devices); i++ {
		res[i] = s.devices[i].ToEntity()
	}

	return res
}
