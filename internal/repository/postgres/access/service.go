package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *AccessRepository) AddRoleToService(ctx context.Context, serviceUid uuid.UUID, role string) error {
	log.Printf("accessRepository.AddRoleToService: serviceUid %s role %s", serviceUid, role)

	stmt := "INSERT INTO services_roles (service_uid, role) VALUES ($1, $2)"

	_, err := r.Exec(ctx, stmt, []interface{}{
		pgtype.UUID{
			Bytes: serviceUid,
			Valid: true,
		},
		role,
	}...)
	if err != nil {
		log.Printf("accessRepository.AddRoleToService: failed to add role %s to service uid %s: %v", role, serviceUid, err)
		return err
	}
	return nil
}

func (r *AccessRepository) DeleteRoleFromService(ctx context.Context, serviceUid uuid.UUID, role string) error {
	log.Printf("accessRepository.DeleteRoleFromService: serviceUid %s role %s", serviceUid, role)

	stmt := "DELETE FROM services_roles WHERE service_uid=$1 AND role=$2"

	_, err := r.Exec(ctx, stmt, []interface{}{
		pgtype.UUID{
			Bytes: serviceUid,
			Valid: true,
		},
		role,
	}...)
	if err != nil {
		log.Printf("accessRepository.DeleteRoleFromService: failed to delete role %s from service uid %s: %v", role, serviceUid, err)
		return err
	}

	return nil
}

func (r *AccessRepository) GetServiceRoles(ctx context.Context, serviceUid uuid.UUID) ([]entity.Role, error) {
	log.Printf("accessRepository.GetServiceRoles: service uid %s", serviceUid)

	stmt := "SELECT sr.role, r.description from services_roles sr INNER JOIN roles r on sr.role=r.role  WHERE sr.service_uid=$1"

	rows, err := r.Query(ctx, stmt, pgtype.UUID{
		Bytes: serviceUid,
		Valid: true,
	})
	if err != nil {
		log.Printf("accessRepository.GetServiceRoles: failed to get service %s roles: %v", serviceUid, err)
		return nil, err
	}

	var roles []entity.Role

	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.Role, &role.Description); err != nil {
			log.Printf("accessRepository.GetUserRoles: failed to scan service %s roles: %v", serviceUid, err)
			return nil, errors.WithStack(err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}
