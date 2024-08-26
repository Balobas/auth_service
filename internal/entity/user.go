package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Uid         uuid.UUID
	Name        string
	Email       string
	Phone       string
	Password    string
	Role        string
	Permissions []UserPermission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserPermission string

const (
	UserPermissionNotVerified = UserPermission("not_verified")
	UserPermissionBase        = UserPermission("base")
)
