package pgEntity

import (
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const usersBlackListTableName = "black_list_users"

type UsersBlackListRow struct {
	UserUid pgtype.UUID
	Reason  string
}

func NewUsersBlackListRow() *UsersBlackListRow {
	return &UsersBlackListRow{}
}

var (
	usersBlackListColumns = []string{
		"user_uid",
		"reason",
	}
)

func (u *UsersBlackListRow) FromBlackListEntity(elem entity.BlackListElement) *UsersBlackListRow {
	u.UserUid = pgtype.UUID{
		Bytes:  elem.Uid,
		Status: pgtype.Present,
	}
	u.Reason = elem.Reason
	return u
}

func (u *UsersBlackListRow) ToBlackListEntity() entity.BlackListElement {
	return entity.BlackListElement{
		Uid:    u.UserUid.Bytes,
		Reason: u.Reason,
	}
}

func (u *UsersBlackListRow) FromBlackListUserEntity(user entity.BlackListUser) *UsersBlackListRow {
	u.UserUid = pgtype.UUID{
		Bytes:  user.User.Uid,
		Status: pgtype.Present,
	}
	u.Reason = user.Reason
	return u
}

func (u *UsersBlackListRow) IdColumnName() string {
	return "user_uid"
}

func (u *UsersBlackListRow) Values() []interface{} {
	return []interface{}{
		u.UserUid,
		u.Reason,
	}
}

func (u *UsersBlackListRow) Columns() []string {
	return usersBlackListColumns
}

func (u *UsersBlackListRow) Table() string {
	return usersBlackListTableName
}

func (u *UsersBlackListRow) GetId() interface{} {
	return u.UserUid
}

func (u *UsersBlackListRow) ScanId(row pgx.Row) error {
	return row.Scan(&u.UserUid)
}

func (u *UsersBlackListRow) Scan(row pgx.Row) error {
	return row.Scan(&u.UserUid, &u.Reason)
}

func (u *UsersBlackListRow) ColumnsForUpdate() []string {
	return []string{"reason"}
}

func (u *UsersBlackListRow) ValuesForUpdate() []interface{} {
	return []interface{}{u.Reason}
}

func (u *UsersBlackListRow) ValuesForScan() []interface{} {
	return []interface{}{
		u.UserUid,
		u.Reason,
	}
}
