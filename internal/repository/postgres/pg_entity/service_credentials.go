package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

type ServiceCredentialsRow struct {
	ServiceUid   pgtype.UUID
	PasswordHash pgtype.Text
}

const serviceCredentialsTableName = "services_credentials"

var (
	servicesCredentialsColumns = []string{
		"service_uid",
		"h_password",
	}
)

func NewServiceCredentialsRow() *ServiceCredentialsRow {
	return &ServiceCredentialsRow{}
}

func (uc *ServiceCredentialsRow) FromEntity(creds entity.ServiceCredentials) *ServiceCredentialsRow {
	uc.ServiceUid = pgtype.UUID{
		Bytes: creds.ServiceUid,
		Valid: true,
	}
	uc.PasswordHash = pgtype.Text{
		String: string(creds.PasswordHash),
		Valid:  true,
	}
	return uc
}

func (uc *ServiceCredentialsRow) ToEntity() entity.ServiceCredentials {
	return entity.ServiceCredentials{
		ServiceUid:   uc.ServiceUid.Bytes,
		PasswordHash: []byte(uc.PasswordHash.String),
	}
}

func (uc *ServiceCredentialsRow) IdColumnName() string {
	return "service_uid"
}

func (uc *ServiceCredentialsRow) Values() []interface{} {
	return []interface{}{
		uc.ServiceUid,
		uc.PasswordHash,
	}
}

func (uc *ServiceCredentialsRow) Columns() []string {
	return servicesCredentialsColumns
}

func (uc *ServiceCredentialsRow) Table() string {
	return serviceCredentialsTableName
}

func (uc *ServiceCredentialsRow) Scan(row pgx.Row) error {
	return row.Scan(
		&uc.ServiceUid,
		&uc.PasswordHash,
	)
}

func (uc *ServiceCredentialsRow) ColumnsForUpdate() []string {
	return []string{
		"h_password",
	}
}

func (uc *ServiceCredentialsRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		uc.PasswordHash,
	}
}

func (uc *ServiceCredentialsRow) ConditionServiceUidEqual(uid uuid.UUID) sq.Eq {
	return sq.Eq{
		"service_uid": pgtype.UUID{
			Bytes: uid,
			Valid: true,
		},
	}
}
