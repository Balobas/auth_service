package pgEntity

import (
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	usersTableName = "users"
)

var (
	usersTableColumns = []string{
		"uid",
		"email",
		"is_verified",
		"created_at",
		"updated_at",
	}
)

type UserRow struct {
	Uid        pgtype.UUID
	Email      string
	IsVerified bool
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

func NewUserRow() *UserRow {
	return &UserRow{}
}

type UserRows struct {
	users []*UserRow
}

func NewUserRows() *UserRows {
	return &UserRows{}
}

func (s *UserRows) ScanAll(rows pgx.Rows) error {
	for rows.Next() {
		newRow := &UserRow{}

		if err := newRow.Scan(rows); err != nil {
			return err
		}
		s.users = append(s.users, newRow)
	}

	return nil
}

func (ur *UserRow) Table() string {
	return usersTableName
}

func (ur *UserRow) FromEntity(user entity.User) *UserRow {
	ur.Uid = pgtype.UUID{
		Bytes: user.Uid,
		Valid: true,
	}
	ur.Email = user.Email

	ur.IsVerified = user.IsVerified

	if user.CreatedAt.Unix() == 0 {
		ur.CreatedAt = pgtype.Timestamp{}
	} else {
		ur.CreatedAt = pgtype.Timestamp{
			Time:  user.CreatedAt.UTC(),
			Valid: true,
		}
	}

	if user.UpdatedAt.Unix() == 0 {
		ur.UpdatedAt = pgtype.Timestamp{}
	} else {
		ur.UpdatedAt = pgtype.Timestamp{
			Time:  user.UpdatedAt.UTC(),
			Valid: true,
		}
	}

	return ur
}

func (ur *UserRow) ToEntity() entity.User {
	return entity.User{
		Uid:        ur.Uid.Bytes,
		Email:      ur.Email,
		IsVerified: ur.IsVerified,
		CreatedAt:  ur.CreatedAt.Time,
		UpdatedAt:  ur.UpdatedAt.Time,
	}
}

func (ur *UserRow) Columns() []string {
	return usersTableColumns
}

func (ur *UserRow) Values() []interface{} {
	return []interface{}{
		ur.Uid,
		ur.Email,
		ur.IsVerified,
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
		&ur.IsVerified,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	)
}

func (ur *UserRow) ValuesForScan() []interface{} {
	return []interface{}{
		&ur.Uid,
		&ur.Email,
		&ur.IsVerified,
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

func (ur *UserRow) ConditionNotSuperAdmin() sq.NotEq {
	return sq.NotEq{
		"email": os.Getenv("SUPER_ADMIN_EMAIL"),
	}
}

func (ur *UserRow) ConditionEmailEqual() sq.Eq {
	return sq.Eq{
		"email": ur.Email,
	}
}

func (s *UserRows) ToEntity() []entity.User {
	if len(s.users) == 0 {
		return nil
	}

	res := make([]entity.User, len(s.users))

	for i := 0; i < len(s.users); i++ {
		res[i] = s.users[i].ToEntity()
	}

	return res
}
