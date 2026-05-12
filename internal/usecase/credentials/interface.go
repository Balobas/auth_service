package useCaseCredentials

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type Config interface {
}

type CredentialsRepository interface {
	CreateCredentials(ctx context.Context, creds entity.UserCredentials) error
	UpdateCredentials(ctx context.Context, creds entity.UserCredentials) error
	GetByUserUid(ctx context.Context, userUid uuid.UUID) (entity.UserCredentials, bool, error)
	DeleteByUserUid(ctx context.Context, userUid uuid.UUID) error

	CreateServiceCredentials(ctx context.Context, creds entity.ServiceCredentials) error 
	GetByServiceUid(ctx context.Context, serviceUid uuid.UUID) (entity.ServiceCredentials, bool, error)
	
}
