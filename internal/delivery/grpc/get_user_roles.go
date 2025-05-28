package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (s *AuthServerGrpc) GetUserRoles(ctx context.Context, req *auth_v1.UUID) (*auth_v1.Roles, error) {
	log.Printf("authServerGrpc.GetUserRoles: user %s", req.GetUid())

	userInfo := userInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetUserRoles: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	roles, err := s.ucAccess.GetUserRoles(ctx, uuid.FromStringOrNil(req.GetUid()))
	if err != nil {
		log.Printf("authServerGrpc.GetUserRoles: failed to get user roles: %v", err)
		return nil, errors.WithStack(err)
	}

	rolesPb := make([]*auth_v1.Role, len(roles))
	for i, role := range roles {
		rolesPb[i] = &auth_v1.Role{Name: role.Role, Description: role.Description}
	}

	return &auth_v1.Roles{Roles: rolesPb}, nil
}
