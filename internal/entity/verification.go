package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Verification struct {
	UserUid   uuid.UUID
	Token     string
	Status    VerificationStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VerificationStatus string

const (
	VerificationStatusCreated VerificationStatus = "created"
	VerificationStatusWaiting VerificationStatus = "waiting"
)
