package pgEntity

import (
	basePgEntity "github.com/balobas/sport_city_common/repository/postgres/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

type UidRow struct {
	Uid pgtype.UUID
}

func (u *UidRow) New() *UidRow {
	return &UidRow{}
}

func (u *UidRow) ToEntity() uuid.UUID {
	return u.Uid.Bytes
}

func (u *UidRow) Scan(row pgx.Row) error {
	return row.Scan(&u.Uid)
}

func NewUUIDRows() *basePgEntity.Rows[*UidRow, uuid.UUID] {
	return &basePgEntity.Rows[*UidRow, uuid.UUID]{}
}
