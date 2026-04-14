package entity

import uuid "github.com/satori/go.uuid"

type TokenInfo struct {
	UserUid    uuid.UUID
	DeviceUid  uuid.UUID
	Email      string
	Roles      []string
	SessionUid uuid.UUID
	ExpiredAt  int64
	IssuedAt   int64
}
