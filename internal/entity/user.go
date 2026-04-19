package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Uid        uuid.UUID
	Email      string
	Roles      []string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserInfo struct {
	UserUid uuid.UUID
	Roles   []string
	Token   string
}
