package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type UserCredentialsRow struct {
	UserUid      pgtype.UUID
	PasswordHash pgtype.Varchar
}

const userCredentialsTableName = "users_credentials"

var (
	userCredentialsColumns = []string{
		"user_uid",
		"h_password",
	}
)

func NewUserCredentialsRow() *UserCredentialsRow {
	return &UserCredentialsRow{}
}

func (uc *UserCredentialsRow) FromEntity(creds entity.UserCredentials) *UserCredentialsRow {
	uc.UserUid = pgtype.UUID{
		Bytes:  creds.UserUid,
		Status: pgtype.Present,
	}
	uc.PasswordHash = pgtype.Varchar{
		String: string(creds.PasswordHash),
		Status: pgtype.Present,
	}
	return uc
}

func (uc *UserCredentialsRow) ToEntity() entity.UserCredentials {
	return entity.UserCredentials{
		UserUid:      uc.UserUid.Bytes,
		PasswordHash: []byte(uc.PasswordHash.String),
	}
}

func (uc *UserCredentialsRow) IdColumnName() string {
	return "user_uid"
}

func (uc *UserCredentialsRow) Values() []interface{} {
	return []interface{}{
		uc.UserUid,
		uc.PasswordHash,
	}
}

func (uc *UserCredentialsRow) Columns() []string {
	return userCredentialsColumns
}

func (uc *UserCredentialsRow) Table() string {
	return userCredentialsTableName
}

func (uc *UserCredentialsRow) Scan(row pgx.Row) error {
	return row.Scan(
		&uc.UserUid,
		&uc.PasswordHash,
	)
}

func (uc *UserCredentialsRow) ColumnsForUpdate() []string {
	return []string{
		"h_password",
	}
}

func (uc *UserCredentialsRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		uc.PasswordHash,
	}
}

func (uc *UserCredentialsRow) ConditionUserUidEqual() sq.Eq {
	return sq.Eq{"user_uid": uc.UserUid}
}
