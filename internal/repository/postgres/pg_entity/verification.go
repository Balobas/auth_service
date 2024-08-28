package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const verificationTableName = "verification"

var verificationTableColumns = []string{
	"user_uid",
	"token",
}

type VerificationRow struct {
	UserUid pgtype.UUID
	Token   string
}

func NewVerificationRow() *VerificationRow {
	return &VerificationRow{}
}

func (v *VerificationRow) FromEntity(verification entity.Verification) *VerificationRow {
	v.UserUid = pgtype.UUID{
		Bytes:  verification.UserUid,
		Status: pgtype.Present,
	}
	v.Token = verification.Token
	return v
}

func (v *VerificationRow) ToEntity() entity.Verification {
	return entity.Verification{
		UserUid: v.UserUid.Bytes,
		Token:   v.Token,
	}
}

func (v *VerificationRow) IdColumnName() string {
	return "user_uid"
}

func (v *VerificationRow) Values() []interface{} {
	return []interface{}{
		v.UserUid,
		v.Token,
	}
}

func (v *VerificationRow) Columns() []string {
	return verificationTableColumns
}

func (v *VerificationRow) Table() string {
	return verificationTableName
}

func (v *VerificationRow) Scan(row pgx.Row) error {
	return row.Scan(&v.UserUid, &v.Token)
}

func (v *VerificationRow) ColumnsForUpdate() []string {
	return nil
}

func (v *VerificationRow) ValuesForUpdate() []interface{} {
	return nil
}

func (v *VerificationRow) ConditionUserUidEqual() sq.Eq {
	return sq.Eq{
		"user_uid": v.UserUid,
	}
}
