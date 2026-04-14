package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Uid            uuid.UUID
	UserUid        uuid.UUID
	DeviceUid      uuid.UUID
	TokensIssuedAt int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
