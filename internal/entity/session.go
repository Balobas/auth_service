package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Uid       uuid.UUID
	UserUid   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
