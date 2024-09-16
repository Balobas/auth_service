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
	"email",
	"token",
	"status",
	"created_at",
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

func (v *VerificationRow) FromEntity(verification entity.Verification) *VerificationRow {
	v.UserUid = pgtype.UUID{
		Bytes:  verification.UserUid,
		Status: pgtype.Present,
	}
	v.Token = verification.Token
	v.Email = verification.Email
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

type VerificationRows struct {
	verifications []*VerificationRow
}

func NewVerificationRows() *VerificationRows {
	return &VerificationRows{}
}

func (s *VerificationRows) ScanAll(rows pgx.Rows) error {
	for rows.Next() {
		newRow := &VerificationRow{}

		if err := newRow.Scan(rows); err != nil {
			return err
		}
		s.verifications = append(s.verifications, newRow)
	}

	return nil
}

func (s *VerificationRows) ToEntities() []entity.Verification {
	if len(s.verifications) == 0 {
		return nil
	}

	res := make([]entity.Verification, len(s.verifications))

	for i := 0; i < len(s.verifications); i++ {
		res[i] = s.verifications[i].ToEntity()
	}

	return res
}
