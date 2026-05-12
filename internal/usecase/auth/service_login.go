package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	common "github.com/balobas/sport_city_common"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) ServiceLogin(ctx context.Context, params entity.ServiceLoginParams) (string, string, error) {
	log.Printf("usecaseAuth.ServiceLogin: uid %s, device %s", params.Uid, params.Device.Uid)

	loginTime := time.Now()

	var access, refresh string
	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		service, err := uc.getServiceWithRolesByUid(ctx, params.Uid)
		if err != nil {
			return err
		}

		device, err := entity.NewSystemDevice(service.Uid, params.Device)
		if err != nil {
			return err
		}

		if err := uc.ucCredentials.ValidateServiceCreds(ctx, service.Uid, params.Password); err != nil {
			return err
		}

		savedSession, isFound, err := uc.sessionsRepo.GetSessionByMaintainerAndDevice(ctx, service.Uid, device.Uid)
		if err != nil {
			return err
		}
		if isFound {
			access, refresh, err = uc.handleServiceLoginWithExistingSession(service, savedSession)
			return err
		}

		if err := uc.ucDevices.HandleLoginFromDevice(ctx, device, loginTime); err != nil {
			return err
		}

		session := entity.NewSystemSession(uuid.NewV4(), service.Uid, device.Uid, loginTime)
		if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
			return err
		}

		tokenInfo := entity.NewSystemTokenInfo(service, device.Uid, session.Uid, loginTime)
		access, refresh, err = uc.jwtManager.NewSystemTokens(tokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
		return err
	}); err != nil {
		return emptyTokensWithError(err)
	}

	return access, refresh, nil
}

func (uc *UseCaseAuth) getServiceWithRolesByUid(ctx context.Context, uid uuid.UUID) (entity.Service, error) {
	service, isFound, err := uc.servicesRepo.GetServiceByUid(ctx, uid)
	if err != nil {
		return entity.Service{}, err
	}
	if !isFound {
		return entity.Service{}, serviceErrors.ErrNotFound
	}
	roles, err := uc.accessRepo.GetServiceRoles(ctx, service.Uid)
	if err != nil {
		return entity.Service{}, err
	}

	rolesStrs := entity.RolesToStrings(roles)
	service = service.WithRoles(rolesStrs)
	return service, nil
}

func (uc *UseCaseAuth) handleServiceLoginWithExistingSession(service entity.Service, session entity.Session) (string, string, error) {
	tokenInfo := entity.SystemTokenInfoFromSession(service, session)
	return uc.jwtManager.NewSystemTokens(tokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
}
