package entity

import uuid "github.com/satori/go.uuid"

type Session struct {
	Uid       uuid.UUID
	UserUid   uuid.UUID
	DeviceUid uuid.UUID
	CreatedAt uuid.UUID
	UpdatedAt uuid.UUID
}
