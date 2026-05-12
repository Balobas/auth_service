package deliveryGrpcAuth

import (
	"context"
	"log"

	deliveryGrpcInterceptors "github.com/balobas/auth_service/internal/delivery/grpc/interceptors"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (s *AuthServerGrpc) GetRoles(ctx context.Context, req *auth_v1.GetRolesParams) (*auth_v1.Roles, error) {
	log.Printf("authServerGrpc.GetRoles: role pattern %s", req.GetRolePattern())

	userInfo := deliveryGrpcInterceptors.UserInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetRoles: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	roles, err := s.ucAccess.GetRoles(ctx, req.GetRolePattern(), req.GetLimit(), req.GetOffset())
	if err != nil {
		log.Printf("authServerGrpc.GetRoles: failed to get roles: %v", err)
		return nil, errors.WithStack(err)
	}

	rolesPb := make([]*auth_v1.Role, len(roles))
	for i, role := range roles {
		rolesPb[i] = &auth_v1.Role{Name: role.Role, Description: role.Description}
	}

	return &auth_v1.Roles{Roles: rolesPb}, nil
}
