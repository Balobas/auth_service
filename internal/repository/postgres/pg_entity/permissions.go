package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5"
	basePgEntity "github.com/balobas/sport_city_common/repository/postgres/entity"
)

const permissionsTableName = "permissions"

type PermissionRow struct {
	Key         string
	Description string
}

var permissionsTableColumns = []string{
	"key",
	"description",
}

func NewPermissionRow() *PermissionRow {
	return &PermissionRow{}
}

func (p *PermissionRow) New() *PermissionRow {
	return &PermissionRow{}
}

func (p *PermissionRow) FromEntity(perm entity.Permission) *PermissionRow {
	p.Key = perm.Key
	p.Description = perm.Description
	return p
}

func (p *PermissionRow) ToEntity() entity.Permission {
	return entity.Permission{
		Key:         p.Key,
		Description: p.Description,
	}
}

func (p *PermissionRow) IdColumnName() string {
	return "key"
}

func (p *PermissionRow) Values() []interface{} {
	return []interface{}{
		p.Key,
		p.Description,
	}
}

func (p *PermissionRow) Columns() []string {
	return permissionsTableColumns
}

func (p *PermissionRow) Table() string {
	return permissionsTableName
}

func (p *PermissionRow) Scan(row pgx.Row) error {
	return row.Scan(&p.Key, &p.Description)
}

func (p *PermissionRow) ColumnsForUpdate() []string {
	return []string{"description"}
}

func (p *PermissionRow) ValuesForUpdate() []interface{} {
	return []interface{}{p.Description}
}

func (p *PermissionRow) ConditionUidEqual() sq.Eq {
	return sq.Eq{
		"key": p.Key,
	}
}

func NewPermissionsRows() *basePgEntity.Rows[*PermissionRow, entity.Permission] {
	return &basePgEntity.Rows[*PermissionRow, entity.Permission]{}
}
