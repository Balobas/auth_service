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
	"status",
	"created_at",
}

type VerificationRow struct {
	UserUid   pgtype.UUID
	Token     string
	Status    string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
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
	v.Status = string(verification.Status)
	if verification.CreatedAt.Unix() == 0 {
		v.CreatedAt = pgtype.Timestamp{
			Status: pgtype.Null,
		}
	} else {
		v.CreatedAt = pgtype.Timestamp{
			Time:   verification.CreatedAt,
			Status: pgtype.Present,
		}
	}
	if verification.UpdatedAt.Unix() == 0 {
		v.UpdatedAt = pgtype.Timestamp{
			Status: pgtype.Null,
		}
	} else {
		v.UpdatedAt = pgtype.Timestamp{
			Time:   verification.UpdatedAt,
			Status: pgtype.Present,
		}
	}
	return v
}

func (v *VerificationRow) ToEntity() entity.Verification {
	return entity.Verification{
		UserUid:   v.UserUid.Bytes,
		Token:     v.Token,
		Status:    entity.VerificationStatus(v.Status),
		CreatedAt: v.CreatedAt.Time,
		UpdatedAt: v.UpdatedAt.Time,
	}
}

func (v *VerificationRow) IdColumnName() string {
	return "user_uid"
}

func (v *VerificationRow) Values() []interface{} {
	return []interface{}{
		v.UserUid,
		v.Token,
		v.Status,
		v.CreatedAt,
		v.UpdatedAt,
	}
}

func (v *VerificationRow) Columns() []string {
	return verificationTableColumns
}

func (v *VerificationRow) Table() string {
	return verificationTableName
}

func (v *VerificationRow) Scan(row pgx.Row) error {
	return row.Scan(&v.UserUid, &v.Token, &v.Status, &v.CreatedAt, &v.UpdatedAt)
}

func (v *VerificationRow) ColumnsForUpdate() []string {
	return []string{
		"status",
		"updated_at",
	}
}

func (v *VerificationRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		v.Status,
		v.UpdatedAt,
	}
}

func (v *VerificationRow) ConditionUserUidEqual() sq.Eq {
	return sq.Eq{
		"user_uid": v.UserUid,
	}
}
