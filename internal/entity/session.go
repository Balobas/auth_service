package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Uid            uuid.UUID
	MaintainerUid  uuid.UUID
	Type           SessionType
	DeviceUid      uuid.UUID
	TokensIssuedAt int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SessionType string

const (
	SessionTypeUser   = SessionType("user_session")
	SessionTypeSystem = SessionType("system_session")
)

func NewUserSession(uid uuid.UUID, userUid uuid.UUID, deviceUid uuid.UUID, loginTime time.Time) Session {
	return Session{
		Uid:            uid,
		MaintainerUid:  userUid,
		Type:           SessionTypeUser,
		DeviceUid:      deviceUid,
		TokensIssuedAt: loginTime.UTC().Unix(),
		CreatedAt:      loginTime,
		UpdatedAt:      loginTime,
	}
}

func NewSystemSession(uid uuid.UUID, serviceUid uuid.UUID, deviceUid uuid.UUID, loginTime time.Time) Session {
	return Session{
		Uid:            uid,
		MaintainerUid:  serviceUid,
		Type:           SessionTypeSystem,
		DeviceUid:      deviceUid,
		TokensIssuedAt: loginTime.UTC().Unix(),
		CreatedAt:      loginTime,
		UpdatedAt:      loginTime,
	}
}
