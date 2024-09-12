package entity

import uuid "github.com/satori/go.uuid"

type TokenInfo struct {
	UserUid     uuid.UUID
	Email       string
	Permissions []string
	Role        string
	SessionUid  uuid.UUID
	ExpiredAt   int64
}
