package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthServerGrpc) GetAdmins(ctx context.Context, _ *emptypb.Empty) (*auth_v1.AdminsResponse, error) {
	log.Printf("authServerGrpc.GetAdmins")

	userInfo := userInfoFromContext(ctx)

	if !entity.HasRole(userInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("authServerGrpc.GetAdmins: caller user %s is not admin. permission denied", userInfo.UserUid)
		return nil, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller user is not admin")
	}

	users, err := s.ucUsers.GetAdmins(ctx)
	if err != nil {
		log.Printf("authServerGrpc.GetAdmins: failed to get admin users")
		return nil, err
	}

	respUsers := make([]*auth_v1.GetUserResponse, len(users))
	for ind, user := range users {
		respUsers[ind] = &auth_v1.GetUserResponse{
			Uid:       user.Uid.String(),
			Email:     user.Email,
			Roles:     user.Roles,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		}
	}

	return &auth_v1.AdminsResponse{
		Users: respUsers,
	}, nil
}
