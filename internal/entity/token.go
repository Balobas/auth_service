package entity

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type TokenInfo struct {
	DeviceUid  uuid.UUID
	SessionUid uuid.UUID
	Roles      []string
	IssuedAt   int64

	Type TokenType

	// User token fields
	UserUid uuid.UUID
	Email   string

	// System token fields
	ServiceUid  uuid.UUID
	ServiceName string
	Domain      string

	ExpiredAt int64
}

type TokenType string

const (
	TokenTypeUser   = TokenType("user")
	TokenTypeSystem = TokenType("system")
)

func (tt TokenType) Validate() error {
	switch tt {
	case TokenTypeUser, TokenTypeSystem:
		return nil
	default:
		return fmt.Errorf("invalid token type")
	}
}

func NewUserTokenInfo(user User, deviceUid uuid.UUID, sessionUid uuid.UUID, issuedTime time.Time) TokenInfo {
	return TokenInfo{
		DeviceUid:  deviceUid,
		SessionUid: sessionUid,
		Roles:      user.Roles,
		IssuedAt:   issuedTime.UTC().Unix(),
		Type:       TokenTypeUser,
		UserUid:    user.Uid,
		Email:      user.Email,
	}
}

func NewSystemTokenInfo(service Service, deviceUid uuid.UUID, sessionUid uuid.UUID, issuedTime time.Time) TokenInfo {
	return TokenInfo{
		ServiceUid:  service.Uid,
		ServiceName: service.Name,
		Domain:      service.Domain,
		DeviceUid:   deviceUid,
		SessionUid:  sessionUid,
		Roles:       service.Roles,
		IssuedAt:    issuedTime.UTC().Unix(),
		Type:        TokenTypeSystem,
	}
}

func UserTokenInfoFromSession(user User, session Session) TokenInfo {
	return TokenInfo{
		UserUid:    session.MaintainerUid,
		DeviceUid:  session.DeviceUid,
		Type:       TokenTypeUser,
		Email:      user.Email,
		Roles:      user.Roles,
		SessionUid: session.Uid,
		IssuedAt:   session.TokensIssuedAt,
	}
}

func SystemTokenInfoFromSession(service Service, session Session) TokenInfo {
	return TokenInfo{
		ServiceUid:  service.Uid,
		ServiceName: service.Name,
		Domain:      service.Domain,
		Type:        TokenTypeSystem,
		Roles:       service.Roles,
		SessionUid:  session.Uid,
		IssuedAt:    session.TokensIssuedAt,
	}
}
