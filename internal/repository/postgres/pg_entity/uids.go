package pgEntity

import (
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
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

func NewUUIDRows() *Rows[*UidRow, uuid.UUID] {
	return &Rows[*UidRow, uuid.UUID]{}
}
