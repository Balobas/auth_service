package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (s *AuthServerGrpc) GetRole(ctx context.Context, req *auth_v1.GetRoleParams) (*auth_v1.Role, error) {
	log.Printf("authServerGrpc.GetRole: role %s", req.GetRole())

	userInfo := userInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetRole: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	role, err := s.ucAccess.GetRole(ctx, req.GetRole())
	if err != nil {
		log.Printf("authServerGrpc.GetRole: failed to get role %s: %v", req.GetRole(), err)
		return nil, errors.WithStack(err)
	}

	return &auth_v1.Role{Name: role.Role, Description: role.Description}, nil
}
