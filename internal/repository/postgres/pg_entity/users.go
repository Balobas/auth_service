package pgEntity

import (
	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const (
	usersTableName = "users"
)

var (
	usersTableColumns = []string{
		"id",
		"name",
		"email",
		"password",
		"password_confirm",
		"role",
		"created_at",
		"updated_at",
	}
)

type UserRow struct {
	Id              int64
	Name            string
	Email           string
	Role            string
	Password        string
	ConfirmPassword string
	CreatedAt       pgtype.Timestamp
	UpdatedAt       pgtype.Timestamp
}

func NewUserRow() *UserRow {
	return &UserRow{}
}

func (ur *UserRow) Table() string {
	return usersTableName
}

func (ur *UserRow) FromEntity(user entity.User) *UserRow {
	ur.Id = user.Id
	ur.Name = user.Name
	ur.Email = user.Email
	ur.Role = user.Role
	ur.Password = user.Password
	ur.ConfirmPassword = user.ConfirmPassword

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
		Id:              ur.Id,
		Name:            ur.Name,
		Email:           ur.Email,
		Password:        ur.Password,
		ConfirmPassword: ur.ConfirmPassword,
		Role:            ur.Role,
		CreatedAt:       ur.CreatedAt.Time,
		UpdatedAt:       ur.UpdatedAt.Time,
	}
}

func (ur *UserRow) Columns() []string {
	columns := make([]string, len(usersTableColumns))
	copy(columns, usersTableColumns)
	return columns
}

func (ur *UserRow) ColumnsWithoutId() []string {
	columns := make([]string, len(usersTableColumns)-1)
	copy(columns, usersTableColumns[1:])
	return columns
}

func (ur *UserRow) Values() []interface{} {
	return []interface{}{
		ur.Id,
		ur.Name,
		ur.Email,
		ur.Password,
		ur.ConfirmPassword,
		ur.Role,
		ur.CreatedAt,
		ur.UpdatedAt,
	}
}

func (ur *UserRow) ValuesWithoutId() []interface{} {
	return ur.Values()[1:]
}

func (ur *UserRow) IdColumnName() string {
	return "id"
}

func (ur *UserRow) ScanId(row pgx.Row) error {
	return row.Scan(&ur.Id)
}

func (ur *UserRow) GetId() interface{} {
	return ur.Id
}

func (ur *UserRow) Scan(row pgx.Row) error {
	return row.Scan(
		&ur.Id,
		&ur.Name,
		&ur.Email,
		&ur.Password,
		&ur.ConfirmPassword,
		&ur.Role,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	)
}

func (ur *UserRow) ColumnsForUpdate() []string {
	return []string{
		"name",
		"email",
		"created_at",
	}
}

func (ur *UserRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		ur.Name,
		ur.Email,
		ur.UpdatedAt,
	}
}
