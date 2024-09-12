package entity

import uuid "github.com/satori/go.uuid"

type UserCredentials struct {
	UserUid      uuid.UUID
	PasswordHash []byte
}
