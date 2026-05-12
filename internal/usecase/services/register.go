package useCaseServices

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
)

func (uc *UcServices) Register(ctx context.Context, req entity.RegisterServiceRequest) error {
	log.Printf("ucServices.Register: uid %s, name %s, domain %s", req.Uid, req.Name, req.Domain)

	if len(req.Password) < uc.cfg.MinPasswordLen() {
		log.Printf("ucServices.Register: password shoud have >= %d symbols", uc.cfg.MinPasswordLen())
		return errors.Wrapf(serviceErrors.ErrBadRequest, "password shoud have >= %d symbols", uc.cfg.MinPasswordLen())
	}
	now := time.Now()

	service, err := entity.NewService(req.Uid, req.Name, req.Domain, now)
	if err != nil {
		return err
	}

	return uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {
		_, isFound, err := uc.servicesRepo.GetServiceByDomainAndName(ctx, req.Domain, req.Name)
		if err != nil {
			log.Printf("ucServices.Register: failed to get service by domain %s and name %s: %v", req.Domain, req.Name, err)
			return err
		}
		if isFound {
			log.Printf("ucServices.Register: service with domain %s and name %s already exist", req.Domain, req.Name)
			return errors.Wrap(serviceErrors.ErrAlreadyExists, "service with domain and name already exists")
		}

		if err := uc.servicesRepo.CreateService(ctx, service); err != nil {
			log.Printf("ucServices.Register: failed to create service (name %s, domain %s): %v", req.Domain, req.Name, err)
			return err
		}

		if err := uc.accessRepo.AddRoleToService(ctx, service.Uid, entity.UserRoleUser); err != nil {
			log.Printf("ucServices.Register: failed to add role user to service (name %s, domain %s): %v", req.Domain, req.Name, err)
			return err
		}

		if err := uc.ucCreds.CreateServiceCreds(ctx, service.Uid, req.Password); err != nil {
			log.Printf("ucServices.Register: failed to create service (name %s, domain %s) credentials: %v", req.Domain, req.Name, err)
			return err
		}

		return nil
	})
}
