package serviceErrors

import "fmt"

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrAlreadyExists = fmt.Errorf("already exists")

	ErrTelegramUserNotFound      = fmt.Errorf("telegram account is not linked to any user")
	ErrTelegramAlreadyLinked     = fmt.Errorf("telegram account is already linked to another user")
	ErrAccountAlreadyHasTelegram = fmt.Errorf("account is already linked to another telegram")
	ErrIdempotentOperation       = fmt.Errorf("idempotent operation")
	ErrBadRequest                = fmt.Errorf("bad request")
	ErrNotAllowedByPermissions   = fmt.Errorf("not allowed by permissions")
	ErrInvalidToken              = fmt.Errorf("invalid token")
	ErrTokenExpired              = fmt.Errorf("token expired")
	ErrTokenNotProvided          = fmt.Errorf("token not provided")
)
