package entity

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserDevice struct {
	Uid       uuid.UUID
	UserUid   uuid.UUID
	Name      string
	Agent     string
	Language  string
	CreatedAt time.Time
}

func (ud UserDevice) Validate() error {
	if uuid.Equal(ud.Uid, uuid.UUID{}) {
		return fmt.Errorf("empty device uid")
	}
	return nil
}

func (ud UserDevice) WithUserUid(userUid uuid.UUID) UserDevice {
	ud.UserUid = userUid
	return ud
}

func (ud UserDevice) Authorize(loginTime time.Time) UserAuthorizedDevice {
	return UserAuthorizedDevice{
		UserDevice:   ud,
		AuthorizedAt: loginTime,
	}
}

type UserAuthorizedDevice struct {
	UserDevice
	AuthorizedAt   time.Time
	UnauthorizedAt *time.Time
}

func (ad UserAuthorizedDevice) Unauthorize(logoutTime time.Time) UserAuthorizedDevice {
	ad.UnauthorizedAt = &logoutTime
	return ad
}
