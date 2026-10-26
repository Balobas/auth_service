package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Uid        uuid.UUID
	Email      string
	TelegramID int64 // 0 = не привязан
	Roles      []string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserInfo struct {
	UserUid     uuid.UUID
	Roles       []string
	Token       string
	DeviceUid   uuid.UUID
	IsSystem    bool
	ServiceUid  uuid.UUID
	Domain      string
	ServiceName string
}
