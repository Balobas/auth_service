package serviceErrors

import "fmt"

var (
	ErrNotFound                = fmt.Errorf("not found")
	ErrAlreadyExists           = fmt.Errorf("already exists")
	ErrIdempotentOperation     = fmt.Errorf("idempotent operation")
	ErrBadRequest              = fmt.Errorf("bad request")
	ErrNotAllowedByPermissions = fmt.Errorf("not allowed by permissions")
	ErrInvalidToken            = fmt.Errorf("invalid token")
	ErrTokenExpired            = fmt.Errorf("token expired")
	ErrTokenNotProvided        = fmt.Errorf("token not provided")
)
