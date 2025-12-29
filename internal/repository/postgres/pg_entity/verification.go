package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	basePgEntity "github.com/balobas/sport_city_common/repository/postgres/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const verificationTableName = "verification"

var verificationTableColumns = []string{
	"user_uid",
	"email",
	"token",
	"status",
	"created_at",
	"updated_at",
}

type VerificationRow struct {
	UserUid   pgtype.UUID
	Email     string
	Token     string
	Status    string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func NewVerificationRow() *VerificationRow {
	return &VerificationRow{}
}

func (v *VerificationRow) New() *VerificationRow {
	return &VerificationRow{}
}

func (v *VerificationRow) FromEntity(verification entity.Verification) *VerificationRow {
	v.UserUid = pgtype.UUID{
		Bytes: verification.UserUid,
		Valid: true,
	}
	v.Token = verification.Token
	v.Email = verification.Email
	v.Status = string(verification.Status)
	if verification.CreatedAt.Unix() == 0 {
		v.CreatedAt = pgtype.Timestamp{}
	} else {
		v.CreatedAt = pgtype.Timestamp{
			Time:  verification.CreatedAt.UTC(),
			Valid: true,
		}
	}
	if verification.UpdatedAt.Unix() == 0 {
		v.UpdatedAt = pgtype.Timestamp{}
	} else {
		v.UpdatedAt = pgtype.Timestamp{
			Time:  verification.UpdatedAt.UTC(),
			Valid: true,
		}
	}
	return v
}

func (v *VerificationRow) ToEntity() entity.Verification {
	return entity.Verification{
		UserUid:   v.UserUid.Bytes,
		Email:     v.Email,
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
		v.Email,
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
	return row.Scan(&v.UserUid, &v.Email, &v.Token, &v.Status, &v.CreatedAt, &v.UpdatedAt)
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

func (v *VerificationRow) ConditionsStatusEqual() sq.Eq {
	return sq.Eq{
		"status": v.Status,
	}
}

func (v *VerificationRow) ConditionTokenEqual() sq.Eq {
	return sq.Eq{
		"token": v.Token,
	}
}

func NewVerificationRows() *basePgEntity.Rows[*VerificationRow, entity.Verification] {
	return &basePgEntity.Rows[*VerificationRow, entity.Verification]{}
}
