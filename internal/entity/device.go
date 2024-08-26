package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserDevice struct {
	Uid         uuid.UUID
	UserUid     uuid.UUID
	Name        string
	OS          string
	ConnectedAt time.Time
}
