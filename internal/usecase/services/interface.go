package useCaseServices

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type (
	Config interface {
		MinPasswordLen() int
	}

	ServicesRepository interface {
		CreateService(ctx context.Context, service entity.Service) error
		GetServiceByDomainAndName(ctx context.Context, domain string, name string) (entity.Service, bool, error)
	}

	AccessRepository interface {
		AddRoleToService(ctx context.Context, serviceUid uuid.UUID, role string) error
		DeleteRoleFromService(ctx context.Context, serviceUid uuid.UUID, role string) error
		GetServiceRoles(ctx context.Context, serviceUid uuid.UUID) ([]entity.Role, error)
	}

	UcCredentials interface {
		ValidateServiceCreds(ctx context.Context, serviceUid uuid.UUID, password string) error
		CreateServiceCreds(ctx context.Context, serviceUid uuid.UUID, password string) error
	}
)
