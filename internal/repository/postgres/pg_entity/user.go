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
		"uid",
		"name",
		"email",
		"phone",
		"password",
		"role",
		"created_at",
		"updated_at",
	}
)

type UserRow struct {
	Uid       pgtype.UUID
	Name      string
	Email     string
	Phone     string
	Password  string
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
	ur.Name = user.Name
	ur.Email = user.Email
	ur.Phone = user.Phone
	ur.Role = user.Role
	ur.Password = user.Password

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
		Name:      ur.Name,
		Email:     ur.Email,
		Phone:     ur.Phone,
		Password:  ur.Password,
		Role:      ur.Role,
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
		ur.Name,
		ur.Email,
		ur.Phone,
		ur.Password,
		ur.Role,
		ur.CreatedAt,
		ur.UpdatedAt,
	}
}

func (ur *UserRow) IdColumnName() string {
	return "uid"
}

func (ur *UserRow) ScanId(row pgx.Row) error {
	return row.Scan(&ur.Uid)
}

func (ur *UserRow) GetId() interface{} {
	return ur.Uid
}

func (ur *UserRow) Scan(row pgx.Row) error {
	return row.Scan(
		&ur.Uid,
		&ur.Name,
		&ur.Email,
		&ur.Phone,
		&ur.Password,
		&ur.Role,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	)
}

func (ur *UserRow) ValuesForScan() []interface{} {
	return []interface{}{
		&ur.Uid,
		&ur.Name,
		&ur.Email,
		&ur.Phone,
		&ur.Password,
		&ur.Role,
		&ur.CreatedAt,
		&ur.UpdatedAt,
	}
}

func (ur *UserRow) ColumnsForUpdate() []string {
	return []string{
		"name",
		"phone",
		"email",
		"updated_at",
	}
}

func (ur *UserRow) ValuesForUpdate() []interface{} {
	return []interface{}{
		ur.Name,
		ur.Phone,
		ur.Email,
		ur.UpdatedAt,
	}
}
