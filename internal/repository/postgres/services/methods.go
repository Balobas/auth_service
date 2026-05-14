package servicesRepository

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *Repository) CreateService(ctx context.Context, service entity.Service) error {
	serviceRow := pgEntity.NewServiceRow().FromEntity(service)
	return r.Create(ctx, serviceRow)
}

func (r *Repository) GetServiceByUid(ctx context.Context, uid uuid.UUID) (entity.Service, bool, error) {
	row := pgEntity.NewServiceRow().FromEntity(entity.Service{Uid: uid})
	if err := r.GetOne(ctx, row, row.ConditionUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Service{}, false, nil
		}
		return entity.Service{}, false, errors.Wrapf(err, "failed to get service by uid %s", uid)
	}
	return row.ToEntity(), true, nil
}

func (r *Repository) GetServiceByDomainAndName(ctx context.Context, domain string, name string) (entity.Service, bool, error) {
	row := pgEntity.NewServiceRow()
	if err := r.GetOne(ctx, row, row.ConditionDomainAndNameEqual(domain, name)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Service{}, false, nil
		}
		return entity.Service{}, false, errors.Wrap(err, "failed to get service by domain and name")
	}
	return row.ToEntity(), true, nil
}
