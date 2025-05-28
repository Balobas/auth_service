package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) UpdateResourcePermission(ctx context.Context, req *auth_v1.ResourcePermissionParams) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.UpdateResourcePermission: resource %s method %s permission %s", req.GetUri(), req.GetMethod(), req.GetPermission())

	userInfo := userInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.UpdateResourcePermission: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	if err := s.ucAccess.UpdateResourcePermission(ctx, req.GetUri(), req.GetMethod(), req.GetPermission()); err != nil {
		log.Printf("authServerGrpc.UpdateResourcePermission: failed to update resource %s method %s permission %s: %v", req.GetUri(), req.GetMethod(), req.GetPermission(), err)
		return nil, errors.WithStack(err)
	}

	return &emptypb.Empty{}, nil
}
