package entity

import uuid "github.com/satori/go.uuid"

type Verification struct {
	UserUid uuid.UUID
	Token   string
}
