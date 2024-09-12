package pgEntity

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const (
	usersTableName = "users"
)

var (
	usersTableColumns = []string{
		"uid",
		"email",
		"role",
		"created_at",
		"updated_at",
	}
)

type UserRow struct {
	Uid       pgtype.UUID
	Email     string
	Role      string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func NewUserRow() *UserRow {
	return &UserRow{}
}

func (ur *UserRow) Table() string {
	return usersTableName
}

func (ur *UserRow) FromEntity(user entity.User) *UserRow {
	ur.Uid = pgtype.UUID{
		Bytes:  user.Uid,
		Status: pgtype.Present,
	}
	ur.Email = user.Email
	ur.Role = string(user.Role)

	if user.CreatedAt.Unix() == 0 {
		ur.CreatedAt = pgtype.Timestamp{
			Status: pgtype.Null,
		}
	} else {
		ur.CreatedAt = pgtype.Timestamp{
			Time:   user.CreatedAt,
			Status: pgtype.Present,
		}
	}

	if user.UpdatedAt.Unix() == 0 {
		ur.UpdatedAt = pgtype.Timestamp{
			Status: pgtype.Null,
		}
	} else {
		ur.UpdatedAt = pgtype.Timestamp{
			Time:   user.UpdatedAt,
			Status: pgtype.Present,
		}
	}

	return ur
}

func (ur *UserRow) ToEntity() entity.User {
	return entity.User{
		Uid:       ur.Uid.Bytes,
		Email:     ur.Email,
		Role:      entity.UserRole(ur.Role),
		CreatedAt: ur.CreatedAt.Time,
		UpdatedAt: ur.UpdatedAt.Time,
	}
}

func (ur *UserRow) Columns() []string {
	return usersTableColumns
}

func (ur *UserRow) Values() []interface{} {
	return []interface{}{
		ur.Uid,
		ur.Email,
		ur.Role,
		ur.CreatedAt,
		ur.UpdatedAt,
	}
}

func (ur *UserRow) IdColumnName() string {
	return "uid"
}

func (ur *UserRow) Scan(row pgx.Row) error {
	return row.Scan(
		&ur.Uid,
		&ur.Email,
		&ur.Role,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	)
}

func (ur *UserRow) ValuesForScan() []interface{} {
	return []interface{}{
		&ur.Uid,
		&ur.Email,
		&ur.Role,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	}
}

func (ur *UserRow) ColumnsForUpdate() []string {
	return []string{
		"email",
		"updated_at",
	}
}

func (ur *UserRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		ur.Email,
		ur.UpdatedAt,
	}
}

func (ur *UserRow) ConditionUserUidEqual() sq.Eq {
	return sq.Eq{
		"uid": ur.Uid,
	}
}

func (ur *UserRow) ConditionEmailEqual() sq.Eq {
	return sq.Eq{
		"email": ur.Email,
	}
}
