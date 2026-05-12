package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	servicesTableName = "services"
)

var (
	servicesTableColumns = []string{
		"uid",
		"name",
		"domain",
		"created_at",
		"updated_at",
	}
)

type ServiceRow struct {
	Uid       pgtype.UUID
	Name      string
	Domain    string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func NewServiceRow() *ServiceRow {
	return &ServiceRow{}
}

func (ur *ServiceRow) Table() string {
	return servicesTableName
}

func (ur *ServiceRow) FromEntity(service entity.Service) *ServiceRow {
	ur.Uid = pgtype.UUID{
		Bytes: service.Uid,
		Valid: true,
	}
	ur.Name = service.Name
	ur.Domain = service.Domain

	if service.CreatedAt.Unix() == 0 {
		ur.CreatedAt = pgtype.Timestamp{}
	} else {
		ur.CreatedAt = pgtype.Timestamp{
			Time:  service.CreatedAt.UTC(),
			Valid: true,
		}
	}

	if service.UpdatedAt.Unix() == 0 {
		ur.UpdatedAt = pgtype.Timestamp{}
	} else {
		ur.UpdatedAt = pgtype.Timestamp{
			Time:  service.UpdatedAt.UTC(),
			Valid: true,
		}
	}

	return ur
}

func (ur *ServiceRow) ToEntity() entity.Service {
	return entity.Service{
		Uid:       ur.Uid.Bytes,
		Name:      ur.Name,
		Domain:    ur.Domain,
		CreatedAt: ur.CreatedAt.Time,
		UpdatedAt: ur.UpdatedAt.Time,
	}
}

func (ur *ServiceRow) Columns() []string {
	return servicesTableColumns
}

func (ur *ServiceRow) Values() []interface{} {
	return []interface{}{
		ur.Uid,
		ur.Name,
		ur.Domain,
		ur.CreatedAt,
		ur.UpdatedAt,
	}
}

func (ur *ServiceRow) IdColumnName() string {
	return "uid"
}

func (ur *ServiceRow) Scan(row pgx.Row) error {
	return row.Scan(
		&ur.Uid,
		&ur.Name,
		&ur.Domain,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	)
}

func (ur *ServiceRow) ValuesForScan() []interface{} {
	return []interface{}{
		&ur.Uid,
		&ur.Name,
		&ur.Domain,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	}
}

func (ur *ServiceRow) ColumnsForUpdate() []string {
	return []string{}
}

func (ur *ServiceRow) ValuesForUpdate() []interface{} {
	return []interface{}{}
}

func (ur *ServiceRow) ConditionUidEqual() sq.Eq {
	return sq.Eq{
		"uid": ur.Uid,
	}
}

func (ur *ServiceRow) ConditionDomainAndNameEqual(domain, name string) sq.Eq {
	return sq.Eq{
		"name":   name,
		"domain": domain,
	}
}
