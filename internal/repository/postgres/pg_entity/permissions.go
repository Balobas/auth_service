package pgEntity

import (
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

const userPermissionsTableName = "user_permissions"

type UserPermissionsRow struct {
	UserUid     pgtype.UUID
	Permissions []string
}

var userPermissionsTableColumns = []string{
	"user_uid",
	"permissions",
}

func NewUserPermissionsRow() *UserPermissionsRow {
	return &UserPermissionsRow{}
}

func (p *UserPermissionsRow) FromEntity(user entity.User) *UserPermissionsRow {
	p.UserUid = pgtype.UUID{
		Bytes:  user.Uid,
		Status: pgtype.Present,
	}

	p.Permissions = make([]string, len(user.Permissions))
	for i := 0; i < len(user.Permissions); i++ {
		p.Permissions[i] = string(user.Permissions[i])
	}

	return p
}

func (p *UserPermissionsRow) ToEntity(user *entity.User) {
	if uuid.Equal(user.Uid, uuid.UUID{}) {
		user.Uid = p.UserUid.Bytes
	}
	user.Permissions = make([]entity.UserPermission, len(p.Permissions))
	for i := 0; i < len(p.Permissions); i++ {
		user.Permissions[i] = entity.UserPermission(p.Permissions[i])
	}
}

func (p *UserPermissionsRow) IdColumnName() string {
	return "user_uid"
}

func (p *UserPermissionsRow) Values() []interface{} {
	return []interface{}{
		p.UserUid,
		pq.Array(p.Permissions),
	}
}

func (p *UserPermissionsRow) Columns() []string {
	return userPermissionsTableColumns
}

func (p *UserPermissionsRow) Table() string {
	return userPermissionsTableName
}

func (p *UserPermissionsRow) GetId() interface{} {
	return p.UserUid
}

func (p *UserPermissionsRow) ScanId(row pgx.Row) error {
	return row.Scan(&p.UserUid)
}

func (p *UserPermissionsRow) Scan(row pgx.Row) error {
	return row.Scan(&p.UserUid, &p.Permissions)
}

func (p *UserPermissionsRow) ColumnsForUpdate() []string {
	return []string{"permissions"}
}

func (p *UserPermissionsRow) ValuesForUpdate() []interface{} {
	return []interface{}{pq.Array(p.Permissions)}
}
